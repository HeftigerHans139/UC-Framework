const announcementFetch = typeof authFetch === 'function' ? authFetch : fetch;

function setAnnouncementStatus(id, message, isError) {
    const el = document.getElementById(id);
    if (!el) {
        return;
    }
    el.textContent = message;
    el.className = 'save-status ' + (isError ? 'err' : 'ok');
}

function updateScheduleUI() {
    const repeatEnabled = document.getElementById('announcement-repeat-enabled')?.checked;
    const intervalGroup = document.getElementById('announcement-interval-group');
    const timeGroup = document.getElementById('announcement-time-group');
    const modeRadios = document.querySelectorAll('input[name="announcement-mode"]');

    if (!repeatEnabled) {
        if (intervalGroup) intervalGroup.style.display = 'none';
        if (timeGroup) timeGroup.style.display = 'none';
        return;
    }

    const mode = Array.from(modeRadios).find((r) => r.checked)?.value || 'interval';
    if (intervalGroup) {
        intervalGroup.style.display = mode === 'interval' ? 'block' : 'none';
    }
    if (timeGroup) {
        timeGroup.style.display = mode === 'time' ? 'block' : 'none';
    }
}

async function loadAnnouncementSettings() {
    try {
        const res = await announcementFetch('/api/announcement/settings');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('announcement.settingsLoadError'));
        }
        const settings = await res.json();

        const messageEl = document.getElementById('announcement-message');
        const repeatEl = document.getElementById('announcement-repeat-enabled');
        const modeRadios = document.querySelectorAll('input[name="announcement-mode"]');
        const intervalEl = document.getElementById('announcement-interval');
        const intervalCountEl = document.getElementById('announcement-repeat-count');
        const timeEl = document.getElementById('announcement-time');

        if (messageEl) {
            messageEl.value = settings.message || '';
        }
        if (repeatEl) {
            repeatEl.checked = !!settings.repeat_enabled;
        }
        if (modeRadios) {
            const mode = settings.schedule_mode || 'interval';
            modeRadios.forEach((r) => {
                r.checked = r.value === mode;
            });
        }
        if (intervalEl) {
            intervalEl.value = Math.max(10, settings.repeat_interval_minutes || 60);
        }
        if (intervalCountEl) {
            intervalCountEl.value = Math.max(1, settings.repeat_interval_count || 1);
        }
        if (timeEl) {
            timeEl.value = settings.repeat_time || '08:00';
        }

        updateScheduleUI();
        setAnnouncementStatus('announcement-settings-status', t('announcement.loaded'), false);
    } catch (err) {
        setAnnouncementStatus('announcement-settings-status', t('announcement.settingsLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function saveAnnouncementSettings() {
    const messageEl = document.getElementById('announcement-message');
    const repeatEl = document.getElementById('announcement-repeat-enabled');
    const modeRadios = document.querySelectorAll('input[name="announcement-mode"]');
    const intervalEl = document.getElementById('announcement-interval');
    const intervalCountEl = document.getElementById('announcement-repeat-count');
    const timeEl = document.getElementById('announcement-time');

    const message = (messageEl?.value || '').trim();
    if (!message) {
        setAnnouncementStatus('announcement-settings-status', t('announcement.messageEmpty'), true);
        return;
    }

    if (message.length > 500) {
        setAnnouncementStatus('announcement-settings-status', t('announcement.messageTooLong'), true);
        return;
    }

    const repeatEnabled = repeatEl?.checked || false;
    const scheduleMode = Array.from(modeRadios).find((r) => r.checked)?.value || 'interval';
    const intervalMinutes = Math.max(10, Number(intervalEl?.value || 60));
    const intervalCount = Math.max(1, Number(intervalCountEl?.value || 1));
    const repeatTime = timeEl?.value || '08:00';

    if (repeatEnabled && scheduleMode === 'time') {
        const timeRegex = /^\d{2}:\d{2}$/;
        if (!timeRegex.test(repeatTime)) {
            setAnnouncementStatus('announcement-settings-status', t('announcement.invalidTime'), true);
            return;
        }
    }

    const payload = {
        message,
        repeat_enabled: repeatEnabled,
        schedule_mode: scheduleMode,
        repeat_interval_minutes: intervalMinutes,
        repeat_interval_count: intervalCount,
        repeat_time: repeatTime
    };

    try {
        const res = await announcementFetch('/api/announcement/settings', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('announcement.saveError'));
        }

        setAnnouncementStatus('announcement-settings-status', t('announcement.saved'), false);
    } catch (err) {
        setAnnouncementStatus('announcement-settings-status', t('announcement.saveError', { error: err && err.message ? err.message : err }), true);
    }
}

async function sendAnnouncementNow() {
    const messageEl = document.getElementById('announcement-message');
    const message = (messageEl?.value || '').trim();

    if (!message) {
        setAnnouncementStatus('announcement-action-status', t('announcement.messageEmpty'), true);
        return;
    }

    try {
        const res = await announcementFetch('/api/announcement/send', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message })
        });

        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('announcement.sendError'));
        }

        setAnnouncementStatus('announcement-action-status', t('announcement.sent'), false);
        await loadAnnouncementStatus();
    } catch (err) {
        setAnnouncementStatus('announcement-action-status', t('announcement.sendError', { error: err && err.message ? err.message : err }), true);
    }
}

async function loadAnnouncementStatus() {
    try {
        const res = await announcementFetch('/api/announcement/status');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('announcement.statusLoadError'));
        }
        const payload = await res.json();
        const status = payload && payload.status ? payload.status : {};

        const lastSentEl = document.getElementById('announcement-last-sent');
        if (lastSentEl) {
            if (status.last_sent_at) {
                const dt = new Date(status.last_sent_at);
                if (!Number.isNaN(dt.getTime())) {
                    lastSentEl.textContent = dt.toLocaleString();
                } else {
                    lastSentEl.textContent = status.last_sent_at;
                }
            } else {
                lastSentEl.textContent = t('announcement.neverSent');
            }
        }
    } catch (err) {
        setAnnouncementStatus('announcement-action-status', t('announcement.statusLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const repeatEl = document.getElementById('announcement-repeat-enabled');
    const modeRadios = document.querySelectorAll('input[name="announcement-mode"]');

    if (repeatEl) {
        repeatEl.addEventListener('change', updateScheduleUI);
    }
    modeRadios.forEach((radio) => {
        radio.addEventListener('change', updateScheduleUI);
    });

    loadAnnouncementSettings();
    loadAnnouncementStatus();
});
