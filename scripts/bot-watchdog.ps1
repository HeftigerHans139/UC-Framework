param(
    [Parameter(Mandatory = $true)]
    [string]$SupervisorScript,

    [Parameter(Mandatory = $true)]
    [string]$BotPath,

    [string]$BotArgs = '',
    [string]$WorkDir = '.',

    [Parameter(Mandatory = $true)]
    [string]$StateFile,

    [Parameter(Mandatory = $true)]
    [string]$PidFile,

    [Parameter(Mandatory = $true)]
    [string]$WatchdogPidFile,

    [string]$LogFile = '',
    [int]$MinIntervalSec = 60,
    [int]$MaxIntervalSec = 120
)

$ErrorActionPreference = 'Continue'

function Ensure-ParentDir([string]$FilePath) {
    $dir = Split-Path -Parent $FilePath
    if ($dir -and -not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}

function Read-State {
    if (-not (Test-Path $StateFile)) {
        return @{ desiredRunning = $false }
    }
    try {
        return (Get-Content $StateFile -Raw | ConvertFrom-Json -AsHashtable)
    } catch {
        return @{ desiredRunning = $false }
    }
}

function Write-Log([string]$msg) {
    if (-not [string]::IsNullOrWhiteSpace($LogFile)) {
        Ensure-ParentDir $LogFile
        Add-Content -Path $LogFile -Value ((Get-Date).ToString('s') + ' ' + $msg)
    }
}

if ($MinIntervalSec -lt 60) { $MinIntervalSec = 60 }
if ($MaxIntervalSec -lt $MinIntervalSec) { $MaxIntervalSec = 120 }

Ensure-ParentDir $WatchdogPidFile
Set-Content -Path $WatchdogPidFile -Value $PID -Encoding UTF8
Write-Log "watchdog started pid=$PID"

while ($true) {
    try {
        $state = Read-State
        $desired = [bool]$state.desiredRunning

        if ($desired) {
            $statusJson = & powershell -NoProfile -ExecutionPolicy Bypass -File $SupervisorScript -Action status -BotPath $BotPath -BotArgs $BotArgs -WorkDir $WorkDir -StateFile $StateFile -PidFile $PidFile -LogFile $LogFile
            $running = $false
            if ($statusJson) {
                try {
                    $status = $statusJson | ConvertFrom-Json
                    $running = [bool]$status.running
                } catch {
                    $running = $false
                }
            }

            if (-not $running) {
                Write-Log "bot not running, attempting restart"
                & powershell -NoProfile -ExecutionPolicy Bypass -File $SupervisorScript -Action start -BotPath $BotPath -BotArgs $BotArgs -WorkDir $WorkDir -StateFile $StateFile -PidFile $PidFile -LogFile $LogFile | Out-Null
                $sleep = Get-Random -Minimum $MinIntervalSec -Maximum ($MaxIntervalSec + 1)
                Start-Sleep -Seconds $sleep
                continue
            }
        }

        Start-Sleep -Seconds 10
    } catch {
        Write-Log ("watchdog loop error: " + $_.Exception.Message)
        $sleepErr = Get-Random -Minimum $MinIntervalSec -Maximum ($MaxIntervalSec + 1)
        Start-Sleep -Seconds $sleepErr
    }
}
