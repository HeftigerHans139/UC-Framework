const frameworkFetch = typeof authFetch === 'function' ? authFetch : fetch;

function setFrameworkStatus(id, msg, isErr) {
    const el = document.getElementById(id);
    if (!el) {
        return;
    }
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
}

function formatFrameworkUpdatedAt(ts) {
    if (!ts) {
        return '-';
    }
    const dt = new Date(ts);
    if (Number.isNaN(dt.getTime())) {
        return ts;
    }
    return dt.toLocaleString();
}

async function loadFrameworkSettings() {
    try {
        const res = await frameworkFetch('/api/framework/settings');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('framework.settingsLoadError'));
        }

        const settings = await res.json();
        const mode = settings && settings.platform_mode ? settings.platform_mode : 'teamspeak';
        const input = document.querySelector(`input[name="framework-platform-mode"][value="${mode}"]`)
            || document.querySelector('input[name="framework-platform-mode"][value="teamspeak"]');
        if (input) {
            input.checked = true;
        }

        setFrameworkStatus('framework-platform-status', t('framework.settingsLoaded'), false);
    } catch (err) {
        setFrameworkStatus('framework-platform-status', t('framework.settingsLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function saveFrameworkSettings() {
    const selected = document.querySelector('input[name="framework-platform-mode"]:checked');
    const payload = {
        platform_mode: selected ? selected.value : 'teamspeak'
    };

    try {
        const res = await frameworkFetch('/api/framework/settings', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('framework.saveError'));
        }
        setFrameworkStatus('framework-platform-status', t('framework.saved'), false);
    } catch (err) {
        setFrameworkStatus('framework-platform-status', t('framework.saveError', { error: err && err.message ? err.message : err }), true);
    }
}

async function loadFrameworkInfo() {
    try {
        const res = await frameworkFetch('/api/framework/info');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('framework.infoLoadError'));
        }

        const payload = await res.json();
        const info = payload && payload.info ? payload.info : {};

        const nameEl = document.getElementById('framework-name');
        const versionEl = document.getElementById('framework-version');
        const updatedEl = document.getElementById('framework-updated-at');
        const latestEl = document.getElementById('framework-latest-version');
        const upToDateEl = document.getElementById('framework-is-latest');

        if (nameEl) {
            nameEl.textContent = info.name || 'UC-Framework';
        }
        if (versionEl) {
            versionEl.textContent = info.version || '-';
        }
        if (updatedEl) {
            updatedEl.textContent = formatFrameworkUpdatedAt(info.updated_at);
        }
        if (latestEl) {
            latestEl.textContent = info.latest_version || '-';
        }
        if (upToDateEl) {
            upToDateEl.textContent = info.is_latest ? t('common.yes') : t('common.no');
            upToDateEl.style.color = info.is_latest ? 'var(--success)' : 'var(--danger)';
        }

        setFrameworkStatus('framework-info-status', t('framework.infoLoaded'), false);
    } catch (err) {
        setFrameworkStatus('framework-info-status', t('framework.infoLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function restartFramework() {
    const confirmed = window.confirm(t('framework.restartConfirm'));
    if (!confirmed) {
        return;
    }

    try {
        const res = await frameworkFetch('/api/framework/restart', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('framework.restartFailed'));
        }

        setFrameworkStatus('framework-restart-status', t('framework.restartScheduled'), false);
        setTimeout(() => {
            window.location.reload();
        }, 2500);
    } catch (err) {
        setFrameworkStatus('framework-restart-status', t('framework.restartError', { error: err && err.message ? err.message : err }), true);
    }
}

async function initFrameworkPage() {
    await Promise.all([
        loadFrameworkInfo(),
        loadFrameworkSettings()
    ]);
}

document.addEventListener('DOMContentLoaded', initFrameworkPage);
