function showLoginStatus(msg, isErr) {
    const el = document.getElementById('login-status');
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
}

async function handleLogin(ev) {
    ev.preventDefault();

    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;

    const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    if (!res.ok) {
        if (res.status === 429) {
            const retryAfter = res.headers.get('Retry-After');
            const msg = retryAfter
                ? t('login.tooManyAttempts', { seconds: retryAfter })
                : t('login.tooManyAttemptsGeneric');
            showLoginStatus(msg, true);
            return;
        }
        showLoginStatus(t('login.failed'), true);
        return;
    }

    const data = await res.json();
    if (data && data.token) {
        setAuthToken(data.token);
        if (data.must_change_password) {
            setForcePasswordChangeRequired(true);
            window.location.href = '/auth-settings.html?force_password_change=1';
            return;
        }
        setForcePasswordChangeRequired(false);
        window.location.href = '/';
        return;
    }

    // Falls Auth deaktiviert ist, Frontend-Guard lokal passieren lassen.
    setAuthToken('__auth_disabled__');
    window.location.href = '/';
}

document.addEventListener('DOMContentLoaded', () => {
    const logoutReason = sessionStorage.getItem('uc_framework_logout_reason');
    if (logoutReason === 'inactivity') {
        sessionStorage.removeItem('uc_framework_logout_reason');
        showLoginStatus(t('login.autoLogout'), false);
    }

    if (getAuthToken()) {
        window.location.href = '/';
        return;
    }
    const form = document.getElementById('login-form');
    form.addEventListener('submit', handleLogin);
});

