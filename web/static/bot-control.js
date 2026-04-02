const botFetch = typeof authFetch === 'function' ? authFetch : fetch;

function setStatus(id, msg, isErr) {
    const el = document.getElementById(id);
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
}

function boolText(v) {
    return v ? t('common.yes') : t('common.no');
}

function setText(id, value) {
    const el = document.getElementById(id);
    if (!el) {
        return;
    }
    const text = (value === null || value === undefined || value === '') ? '-' : String(value);
    el.textContent = text;
    el.title = text;
}

async function loadBotStatus() {
    const res = await botFetch('/api/bot/status');
    if (!res.ok) {
        throw new Error(await res.text());
    }
    const s = await res.json();
    document.getElementById('bot-running').textContent = boolText(!!s.running);
    document.getElementById('bot-pid').textContent = s.pid || '-';
    document.getElementById('bot-desired').textContent = boolText(!!s.desired_running);
    document.getElementById('bot-last-action').textContent = s.last_action || '-';
    document.getElementById('bot-last-error').textContent = s.last_error || '-';
}

async function loadWatchdogStatus() {
    const res = await botFetch('/api/bot/watchdog/status');
    if (!res.ok) {
        throw new Error(await res.text());
    }
    const s = await res.json();
    document.getElementById('watchdog-running').textContent = boolText(!!s.running);
    document.getElementById('watchdog-pid').textContent = s.pid || '-';
}

async function loadSystemInfo() {
    const res = await botFetch('/api/bot/system');
    if (!res.ok) {
        throw new Error(await res.text());
    }
    const s = await res.json();
    setText('sys-os', s.os);
    setText('sys-arch', s.arch);
    setText('sys-workdir', s.working_dir);
    setText('sys-bot', s.bot_executable);
    setText('sys-supervisor', s.supervisor_script);
    setText('sys-watchdog', s.watchdog_script);
    setText('sys-logfile', s.log_file);
}

async function refreshLogs() {
    const lineSelect = document.getElementById('log-lines');
    const levelSelect = document.getElementById('log-level');
    const lines = lineSelect ? lineSelect.value : '100';
    const level = levelSelect ? levelSelect.value : 'all';
    const res = await botFetch(`/api/bot/logs?lines=${encodeURIComponent(lines)}&level=${encodeURIComponent(level)}`);
    if (!res.ok) {
        throw new Error(await res.text());
    }
    const s = await res.json();
    const out = document.getElementById('bot-log-output');
    out.textContent = s.content && s.content.trim() !== '' ? s.content : t('bot.noLogLines');
    setStatus('log-status', t('bot.logsLoaded', { lines: s.lines, level: t(`log.level.${s.level}`) }), false);
}

async function botAction(action) {
    const res = await botFetch('/api/bot/action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action })
    });
    if (!res.ok) {
        setStatus('bot-action-status', `${t('common.error')}: ${await res.text()}`, true);
        return;
    }
    setStatus('bot-action-status', t('bot.actionExecuted', { action }), false);
    await refreshBotControl();
}

async function watchdogAction(action) {
    const res = await botFetch('/api/bot/watchdog/action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action })
    });
    if (!res.ok) {
        setStatus('watchdog-action-status', `${t('common.error')}: ${await res.text()}`, true);
        return;
    }
    setStatus('watchdog-action-status', t('bot.watchdogActionExecuted', { action }), false);
    await refreshBotControl();
}

async function refreshBotControl() {
    try {
        await loadBotStatus();
        await loadWatchdogStatus();
        await loadSystemInfo();
    } catch (e) {
        setStatus('bot-action-status', t('bot.statusLoadError', { error: e.message }), true);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    refreshBotControl();
    refreshLogs().catch((e) => setStatus('log-status', t('bot.logsLoadError', { error: e.message }), true));
    const select = document.getElementById('log-lines');
    if (select) {
        select.addEventListener('change', () => {
            refreshLogs().catch((e) => setStatus('log-status', t('bot.logsLoadError', { error: e.message }), true));
        });
    }
    const levelSelect = document.getElementById('log-level');
    if (levelSelect) {
        levelSelect.addEventListener('change', () => {
            refreshLogs().catch((e) => setStatus('log-status', t('bot.logsLoadError', { error: e.message }), true));
        });
    }
    window.setInterval(refreshBotControl, 10000);
    window.setInterval(() => {
        refreshLogs().catch((e) => setStatus('log-status', t('bot.logsLoadError', { error: e.message }), true));
    }, 10000);
});
