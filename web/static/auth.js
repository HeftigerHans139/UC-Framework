const AUTH_TOKEN_KEY = 'uc_framework_auth_token';
const AUTH_LAST_ACTIVITY_KEY = 'uc_framework_last_activity';
const AUTH_FORCE_PASSWORD_CHANGE_KEY = 'uc_framework_force_password_change';
const INACTIVITY_TIMEOUT_MS = 30 * 60 * 1000;
let inactivityWatcherStarted = false;

function isAuthDisabledToken(token) {
    return token === '__auth_disabled__';
}

function getAuthToken() {
    return sessionStorage.getItem(AUTH_TOKEN_KEY) || '';
}

function setAuthToken(token) {
    sessionStorage.setItem(AUTH_TOKEN_KEY, token);
    touchActivity();
}

function setForcePasswordChangeRequired(required) {
    if (required) {
        sessionStorage.setItem(AUTH_FORCE_PASSWORD_CHANGE_KEY, '1');
    } else {
        sessionStorage.removeItem(AUTH_FORCE_PASSWORD_CHANGE_KEY);
    }
}

function isForcePasswordChangeRequired() {
    return sessionStorage.getItem(AUTH_FORCE_PASSWORD_CHANGE_KEY) === '1';
}

function clearAuthToken() {
    sessionStorage.removeItem(AUTH_TOKEN_KEY);
    sessionStorage.removeItem(AUTH_LAST_ACTIVITY_KEY);
    sessionStorage.removeItem(AUTH_FORCE_PASSWORD_CHANGE_KEY);
}

function touchActivity() {
    sessionStorage.setItem(AUTH_LAST_ACTIVITY_KEY, String(Date.now()));
}

function getLastActivity() {
    const raw = sessionStorage.getItem(AUTH_LAST_ACTIVITY_KEY);
    const value = Number(raw);
    if (!Number.isFinite(value) || value <= 0) {
        return 0;
    }
    return value;
}

function bindActivityEvents() {
    const onActivity = () => {
        const token = getAuthToken();
        if (!token || isAuthDisabledToken(token)) return;
        touchActivity();
    };

    ['click', 'mousemove', 'keydown', 'scroll', 'touchstart'].forEach((eventName) => {
        window.addEventListener(eventName, onActivity, { passive: true });
    });
}

function startInactivityWatcher() {
    if (inactivityWatcherStarted) {
        return;
    }
    inactivityWatcherStarted = true;

    bindActivityEvents();

    window.setInterval(() => {
        const token = getAuthToken();
        if (!token || isAuthDisabledToken(token)) {
            return;
        }
        const last = getLastActivity();
        if (last > 0 && Date.now() - last > INACTIVITY_TIMEOUT_MS) {
            logout(true);
        }
    }, 60 * 1000);
}

function requireAuth() {
    const token = getAuthToken();
    if (!token) {
        window.location.href = '/login.html';
        return;
    }

    const onAuthSettings = window.location.pathname.endsWith('/auth-settings.html') || window.location.pathname === '/auth-settings.html';
    if (isForcePasswordChangeRequired() && !onAuthSettings) {
        window.location.href = '/auth-settings.html?force_password_change=1';
        return;
    }

    if (!isAuthDisabledToken(token)) {
        if (!getLastActivity()) {
            touchActivity();
        }
        startInactivityWatcher();
    }
}

async function authFetch(url, options = {}) {
    const token = getAuthToken();
    const opts = { ...options };
    opts.headers = { ...(options.headers || {}) };

    if (token) {
        opts.headers.Authorization = `Bearer ${token}`;
    }

    const res = await fetch(url, opts);
    if (res.status === 401) {
        clearAuthToken();
        if (!window.location.pathname.endsWith('/login.html') && window.location.pathname !== '/login.html') {
            window.location.href = '/login.html';
        }
        return res;
    }
    if (res.status === 403) {
        let body = '';
        try {
            body = (await res.clone().text()).toLowerCase();
        } catch (_) {
            body = '';
        }
        if (body.includes('password change required')) {
            setForcePasswordChangeRequired(true);
            const onAuthSettings = window.location.pathname.endsWith('/auth-settings.html') || window.location.pathname === '/auth-settings.html';
            if (!onAuthSettings) {
                window.location.href = '/auth-settings.html?force_password_change=1';
            }
        }
    }
    if (res.ok && token && !isAuthDisabledToken(token)) {
        touchActivity();
    }
    return res;
}

async function logout(isInactivity = false) {
    try {
        await authFetch('/api/logout', { method: 'POST' });
    } catch (_) {
        // ignore
    }
    clearAuthToken();
    if (isInactivity) {
        sessionStorage.setItem('uc_framework_logout_reason', 'inactivity');
    }
    window.location.href = '/login.html';
}

