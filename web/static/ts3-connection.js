const ts3ConnFetch = typeof authFetch === 'function' ? authFetch : fetch;

function setTS3ConnStatus(message, isError = false) {
    const el = document.getElementById('ts3conn-status');
    if (!el) {
        return;
    }
    el.textContent = message;
    el.className = 'save-status ' + (isError ? 'err' : 'ok');
}

function formatTS3ConnectionError(status) {
    if (!status || !status.last_error) {
        return '-';
    }
    if (status.bot_running === false || status.last_error === 'bot disabled') {
        return t('ts3conn.botDisabledInfo');
    }
    return status.last_error;
}

function applyTS3ConnectionStatus(status) {
    const botRunning = status.bot_running !== false;
    const connected = !!status.connected;
    const connectedEl = document.getElementById('ts3conn-connected');
    if (connectedEl) {
        if (!botRunning) {
            connectedEl.textContent = t('ts3conn.botDisabled');
            connectedEl.style.color = 'var(--danger)';
        } else {
            connectedEl.textContent = connected ? t('common.yes') : t('common.no');
            connectedEl.style.color = connected ? 'var(--success)' : 'var(--danger)';
        }
    }

    const hostEl = document.getElementById('ts3conn-host');
    if (hostEl) {
        hostEl.textContent = status.host || '-';
    }

    const portEl = document.getElementById('ts3conn-port');
    if (portEl) {
        portEl.textContent = String(status.port || '-');
    }

    const checkEl = document.getElementById('ts3conn-check');
    if (checkEl) {
        checkEl.textContent = status.last_check_at || '-';
    }

    const errEl = document.getElementById('ts3conn-error');
    if (errEl) {
        errEl.textContent = formatTS3ConnectionError(status);
    }
}

async function loadTS3ConnectionStatus() {
    try {
        const res = await ts3ConnFetch('/api/ts3/connection/status');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || 'status failed');
        }
        const payload = await res.json();
        applyTS3ConnectionStatus(payload.status || {});
        setTS3ConnStatus(t('ts3conn.statusLoaded'), false);
    } catch (err) {
        setTS3ConnStatus(t('ts3conn.statusLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function runTS3ConnectionTest() {
    try {
        const res = await ts3ConnFetch('/api/ts3/connection/test', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        const payload = await res.json();
        applyTS3ConnectionStatus(payload.status || {});

        if (!payload.ok) {
			if (payload.status && payload.status.bot_running === false) {
				setTS3ConnStatus(t('ts3conn.botDisabledInfo'), true);
				return;
			}
            const err = payload.error || t('ts3conn.testFailed');
            setTS3ConnStatus(t('ts3conn.testError', { error: err }), true);
            return;
        }

        setTS3ConnStatus(t('ts3conn.testOk'), false);
    } catch (err) {
        setTS3ConnStatus(t('ts3conn.testError', { error: err && err.message ? err.message : err }), true);
    }
}

document.addEventListener('DOMContentLoaded', loadTS3ConnectionStatus);
