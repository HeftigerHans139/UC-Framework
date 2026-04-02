[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseApprovedVerbs', '', Justification = 'False positive in editor diagnostics; no unapproved function verb remains in this script.')]
[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidAssignmentToAutomaticVariable', '', Justification = 'False positive in editor diagnostics; script does not assign to automatic variable PID.')]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('start','stop','restart','status')]
    [string]$Action,

    [Parameter(Mandatory = $true)]
    [string]$BotPath,

    [string]$BotArgs = '',
    [string]$WorkDir = '.',

    [Parameter(Mandatory = $true)]
    [string]$StateFile,

    [Parameter(Mandatory = $true)]
    [Alias('PidFile', 'ProcessIdFile')]
    [string]$ProcessFile,

    [string]$LogFile = ''
)

$ErrorActionPreference = 'Stop'

$ensureParentDirectory = {
    param([string]$FilePath)
    $dir = Split-Path -Parent $FilePath
    if ($dir -and -not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}

function Read-State {
    if (-not (Test-Path $StateFile)) {
        return @{ desiredRunning = $false; lastAction = ''; lastError = ''; lastUpdate = '' }
    }
    try {
        return (Get-Content $StateFile -Raw | ConvertFrom-Json -AsHashtable)
    } catch {
        return @{ desiredRunning = $false; lastAction = ''; lastError = ''; lastUpdate = '' }
    }
}

function Write-State([hashtable]$state) {
    & $ensureParentDirectory $StateFile
    $state.lastUpdate = (Get-Date).ToString('s')
    ($state | ConvertTo-Json -Depth 8) | Set-Content -Path $StateFile -Encoding UTF8
}

function Get-RunningProcess {
    if (-not (Test-Path $ProcessFile)) { return $null }

    $processIdText = (Get-Content $ProcessFile -Raw).Trim()
    if (-not $processIdText) {
        Remove-Item $ProcessFile -Force -ErrorAction SilentlyContinue
        return $null
    }

    if ($processIdText -notmatch '^\d+$') {
        Remove-Item $ProcessFile -Force -ErrorAction SilentlyContinue
        return $null
    }
    $processIdValue = [int]$processIdText

    $proc = Get-Process -Id $processIdValue -ErrorAction SilentlyContinue
    if (-not $proc) {
        Remove-Item $ProcessFile -Force -ErrorAction SilentlyContinue
        return $null
    }
    return $proc
}

function Resolve-BotPath {
    $resolved = if ([System.IO.Path]::IsPathRooted($BotPath)) {
        $BotPath
    } else {
        Join-Path $WorkDir $BotPath
    }

    if (Test-Path $resolved) {
        return $resolved
    }

    # Windows fallback: allow config values without .exe extension.
    if ($env:OS -eq 'Windows_NT' -and -not $resolved.EndsWith('.exe', [System.StringComparison]::OrdinalIgnoreCase)) {
        $resolvedExe = "$resolved.exe"
        if (Test-Path $resolvedExe) {
            return $resolvedExe
        }
    }

    throw "Bot executable not found: $resolved"
}

function Split-Args([string]$argsText) {
    if ([string]::IsNullOrWhiteSpace($argsText)) { return @() }
    return $argsText -split '\s+'
}

function Start-BotProcess([string]$resolvedBotPath, [string[]]$argList) {
    $previousLogFile = [System.Environment]::GetEnvironmentVariable('UC_FRAMEWORK_LOG_FILE', 'Process')
    if ($LogFile) {
        [System.Environment]::SetEnvironmentVariable('UC_FRAMEWORK_LOG_FILE', $LogFile, 'Process')
    }
    try {
        if (@($argList).Count -gt 0) {
            return Start-Process -FilePath $resolvedBotPath -ArgumentList $argList -WorkingDirectory $WorkDir -WindowStyle Hidden -PassThru
        }

        return Start-Process -FilePath $resolvedBotPath -WorkingDirectory $WorkDir -WindowStyle Hidden -PassThru
    } finally {
        [System.Environment]::SetEnvironmentVariable('UC_FRAMEWORK_LOG_FILE', $previousLogFile, 'Process')
    }
}

function Out-Json([hashtable]$obj) {
    ($obj | ConvertTo-Json -Compress)
}

& $ensureParentDirectory $ProcessFile
if ($LogFile) { & $ensureParentDirectory $LogFile }

$state = Read-State

switch ($Action) {
    'status' {
        $proc = Get-RunningProcess
        Out-Json @{
            ok = $true
            running = [bool]$proc
            'pid' = if ($proc) { $proc.Id } else { 0 }
            desired_running = [bool]$state.desiredRunning
            last_action = $state.lastAction
            last_error = $state.lastError
        }
        exit 0
    }

    'start' {
        $state.desiredRunning = $true
        $state.lastAction = 'start'

        $proc = Get-RunningProcess
        if ($proc) {
            Write-State $state
            Out-Json @{ ok = $true; running = $true; 'pid' = $proc.Id; message = 'already running'; desired_running = $true }
            exit 0
        }

        try {
            $resolvedBotPath = Resolve-BotPath
            $argList = Split-Args $BotArgs
            $started = Start-BotProcess -resolvedBotPath $resolvedBotPath -argList $argList

            # Catch fast-crash scenarios so API does not report a false successful start.
            Start-Sleep -Milliseconds 800
            $stillRunning = Get-Process -Id $started.Id -ErrorAction SilentlyContinue
            if (-not $stillRunning) {
                throw "Bot process exited immediately after start."
            }

            Set-Content -Path $ProcessFile -Value $started.Id -Encoding UTF8
            $state.lastError = ''
            Write-State $state
            Out-Json @{ ok = $true; running = $true; 'pid' = $started.Id; message = 'started'; desired_running = $true }
            exit 0
        } catch {
            $state.lastError = $_.Exception.Message
            Write-State $state
            Out-Json @{ ok = $false; running = $false; 'pid' = 0; message = 'start failed'; error = $_.Exception.Message; desired_running = $true }
            exit 1
        }
    }

    'stop' {
        $state.desiredRunning = $false
        $state.lastAction = 'stop'

        $proc = Get-RunningProcess
        if ($proc) {
            try {
                Stop-Process -Id $proc.Id -Force -ErrorAction Stop
            } catch {
                $state.lastError = $_.Exception.Message
                Write-State $state
                Out-Json @{ ok = $false; running = $true; 'pid' = $proc.Id; message = 'stop failed'; error = $_.Exception.Message; desired_running = $false }
                exit 1
            }
        }

        Remove-Item $ProcessFile -Force -ErrorAction SilentlyContinue
        $state.lastError = ''
        Write-State $state
        Out-Json @{ ok = $true; running = $false; 'pid' = 0; message = 'stopped'; desired_running = $false }
        exit 0
    }

    'restart' {
        $state.desiredRunning = $true
        $state.lastAction = 'restart'

        $proc = Get-RunningProcess
        if ($proc) {
            try {
                Stop-Process -Id $proc.Id -Force -ErrorAction Stop
            } catch {
                $state.lastError = $_.Exception.Message
                Write-State $state
                Out-Json @{ ok = $false; running = $true; 'pid' = $proc.Id; message = 'restart stop failed'; error = $_.Exception.Message; desired_running = $true }
                exit 1
            }
            Start-Sleep -Milliseconds 400
        }

        try {
            $resolvedBotPath = Resolve-BotPath
            $argList = Split-Args $BotArgs
            $started = Start-BotProcess -resolvedBotPath $resolvedBotPath -argList $argList

            Start-Sleep -Milliseconds 800
            $stillRunning = Get-Process -Id $started.Id -ErrorAction SilentlyContinue
            if (-not $stillRunning) {
                throw "Bot process exited immediately after restart."
            }

            Set-Content -Path $ProcessFile -Value $started.Id -Encoding UTF8
            $state.lastError = ''
            Write-State $state
            Out-Json @{ ok = $true; running = $true; 'pid' = $started.Id; message = 'restarted'; desired_running = $true }
            exit 0
        } catch {
            $state.lastError = $_.Exception.Message
            Write-State $state
            Out-Json @{ ok = $false; running = $false; 'pid' = 0; message = 'restart start failed'; error = $_.Exception.Message; desired_running = $true }
            exit 1
        }
    }
}
