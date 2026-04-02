function setStatus(id, msg, isErr) {
    const el = document.getElementById(id);
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
}

function showForceHintIfNeeded() {
    const hint = document.getElementById('force-password-hint');
    if (!hint) {
        return;
    }
    const forced = typeof isForcePasswordChangeRequired === 'function' && isForcePasswordChangeRequired();
    hint.style.display = forced ? 'block' : 'none';
}

function authModeLabel(mode) {
    switch (mode) {
    case 'none':
        return t('auth.modeNone');
    case 'ranksystem':
        return t('auth.modeRanksystem');
    case 'local_ranksystem':
        return t('auth.modeLocalRanksystem');
    case 'local':
    default:
        return t('auth.modeLocal');
    }
}

function applyAuthModeUI(mode, ranksystemConfigured = true) {
    const normalizedMode = mode || 'local';
    const noLoginMode = normalizedMode === 'none';
    const localVisible = normalizedMode === 'local' || normalizedMode === 'local_ranksystem';
    const ranksystemVisible = normalizedMode === 'ranksystem' || normalizedMode === 'local_ranksystem';

    const localCard = document.getElementById('local-password-card');
    const localUsernameCard = document.getElementById('local-username-card');
    const healthCard = document.getElementById('auth-health-card');
    const disabledWarning = document.getElementById('auth-disabled-warning');
    const modeSelect = document.getElementById('auth-login-mode');
    const activeMode = document.getElementById('auth-active-mode');

    if (localCard) {
        localCard.style.display = localVisible ? 'block' : 'none';
    }
    if (localUsernameCard) {
        localUsernameCard.style.display = localVisible ? 'block' : 'none';
    }
    if (healthCard) {
        healthCard.style.display = (ranksystemVisible && ranksystemConfigured) ? 'block' : 'none';
    }
    if (disabledWarning) {
        disabledWarning.style.display = noLoginMode ? 'block' : 'none';
    }
    if (modeSelect) {
        modeSelect.value = normalizedMode;
    }
    if (activeMode) {
        activeMode.textContent = authModeLabel(normalizedMode);
    }

    if (!(ranksystemVisible && ranksystemConfigured)) {
        const healthStatus = document.getElementById('health-status');
        if (healthStatus) {
            healthStatus.textContent = '';
            healthStatus.className = 'save-status';
        }
    }
}

async function loadAuthMode() {
    const res = await authFetch('/api/auth/mode');
    if (!res.ok) {
        setStatus('health-status', t('auth.modeLoadError'), true);
        return;
    }
    const mode = await res.json();
    document.getElementById('auth-enabled').textContent = mode.enabled ? t('common.yes') : t('common.no');
    document.getElementById('auth-provider').textContent = mode.provider || '-';
    document.getElementById('auth-username').textContent = mode.username || '-';
    applyAuthModeUI(mode.mode, !!mode.ranksystem_configured);

    if (mode.force_password_change && typeof setForcePasswordChangeRequired === 'function') {
        setForcePasswordChangeRequired(true);
    }

    showForceHintIfNeeded();
}

async function saveAuthMode() {
    const modeSelect = document.getElementById('auth-login-mode');
    const selectedMode = modeSelect ? modeSelect.value : 'local';

    const res = await authFetch('/api/auth/mode', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mode: selectedMode })
    });

    if (!res.ok) {
        setStatus('mode-status', `${t('common.error')}: ${await res.text()}`, true);
        return;
    }

    const mode = await res.json();
    document.getElementById('auth-enabled').textContent = mode.enabled ? t('common.yes') : t('common.no');
    document.getElementById('auth-provider').textContent = mode.provider || '-';
    applyAuthModeUI(mode.mode, !!mode.ranksystem_configured);
    setStatus('mode-status', t('auth.modeSaved'), false);
}

async function changeUsername() {
    const currentPassword = document.getElementById('current-password-username').value;
    const newUsername = document.getElementById('new-username').value.trim();

    if (newUsername.length < 3) {
        setStatus('username-status', t('auth.usernameTooShort'), true);
        return;
    }

    const res = await authFetch('/api/auth/username', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            current_password: currentPassword,
            new_username: newUsername
        })
    });

    if (!res.ok) {
        setStatus('username-status', `${t('common.error')}: ${await res.text()}`, true);
        return;
    }

    const data = await res.json();
    document.getElementById('current-password-username').value = '';
    document.getElementById('new-username').value = '';
    document.getElementById('auth-username').textContent = (data && data.username) ? data.username : newUsername;
    setStatus('username-status', t('auth.usernameChanged'), false);
}

async function checkAuthHealth() {
    const res = await authFetch('/api/auth/health');
    if (!res.ok) {
        setStatus('health-status', t('auth.healthFailed'), true);
        return;
    }
    const data = await res.json();
    const msg = `${data.message}${data.status_code ? ' (HTTP ' + data.status_code + ')' : ''}`;
    setStatus('health-status', msg, !data.healthy);
}

async function changePassword() {
    const currentPassword = document.getElementById('current-password').value;
    const newPassword = document.getElementById('new-password').value;

    if (newPassword.length < 8) {
        setStatus('password-status', t('auth.passwordTooShort'), true);
        return;
    }

    const res = await authFetch('/api/auth/password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            current_password: currentPassword,
            new_password: newPassword
        })
    });

    if (!res.ok) {
        setStatus('password-status', `${t('common.error')}: ${await res.text()}`, true);
        return;
    }

    document.getElementById('current-password').value = '';
    document.getElementById('new-password').value = '';
    if (typeof setForcePasswordChangeRequired === 'function') {
        setForcePasswordChangeRequired(false);
    }
    showForceHintIfNeeded();
    setStatus('password-status', t('auth.passwordChanged'), false);

    const returnTo = new URLSearchParams(window.location.search).get('return_to');
    if (returnTo) {
        window.location.href = returnTo;
        return;
    }
    window.location.href = '/';
}

document.addEventListener('DOMContentLoaded', async () => {
    try {
        await loadAuthMode();
    } catch (_) {
        setStatus('health-status', t('auth.settingsLoadError'), true);
    }
});
