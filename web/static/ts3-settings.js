const ts3Fetch = typeof authFetch === 'function' ? authFetch : fetch;

let ts3ChannelsCache = [];

function setTS3Status(msg, isErr = false) {
    const el = document.getElementById('ts3-settings-status');
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
}

async function loadTS3Settings() {
    try {
        const res = await ts3Fetch('/api/settings/ts3');
        if (!res.ok) {
            setTS3Status(t('ts3.loadError'), true);
            return;
        }
        const data = await res.json();
        document.getElementById('ts3-host').value = data.host || '';
        document.getElementById('ts3-query-port').value = data.query_port || 10011;
        document.getElementById('ts3-voice-port').value = data.voice_port || 9987;
        document.getElementById('ts3-query-user').value = data.query_username || '';
        document.getElementById('ts3-query-pass').value = data.query_password || '';
        document.getElementById('ts3-bot-nickname').value = data.bot_nickname || '';
        document.getElementById('ts3-default-channel').value = data.default_channel || '';
        document.getElementById('ts3-query-slowmode').value = data.query_slowmode_ms || 250;
    } catch (_) {
        setTS3Status(t('ts3.loadException'), true);
    }
}

function renderTS3ChannelList(channels) {
    const list = document.getElementById('ts3-channel-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(channels) || channels.length === 0) {
        list.innerHTML = `<div class="channel-empty">${t('ts3.noChannels')}</div>`;
        return;
    }

    list.innerHTML = channels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('ts3.unnamedChannel');
        return `<button type="button" class="channel-item" data-channel-id="${id}" data-channel-name="${name.replace(/"/g, '&quot;')}">${name}<small>#${id}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const selectedId = btn.getAttribute('data-channel-id') || '';
            const selectedName = btn.getAttribute('data-channel-name') || '';
            const target = document.getElementById('ts3-default-channel');
            if (target) {
                target.value = selectedId ? `${selectedName} (${selectedId})` : selectedName;
            }
            closeTS3ChannelPicker();
        });
    });
}

async function openTS3ChannelPicker() {
    const modal = document.getElementById('ts3-channel-modal');
    if (!modal) {
        return;
    }

    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (ts3ChannelsCache.length > 0) {
        renderTS3ChannelList(ts3ChannelsCache);
        return;
    }

    const list = document.getElementById('ts3-channel-list');
    if (list) {
        list.innerHTML = `<div class="channel-empty">${t('common.loading')}</div>`;
    }

    try {
        const res = await ts3Fetch('/api/settings/ts3/channels');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('ts3.channelLoadError'));
        }
        const payload = await res.json();
        ts3ChannelsCache = Array.isArray(payload.channels) ? payload.channels : [];
        renderTS3ChannelList(ts3ChannelsCache);
    } catch (err) {
        if (list) {
            const msg = err && err.message ? err.message : t('ts3.channelLoadError');
            list.innerHTML = `<div class="channel-empty">${msg}</div>`;
        }
    }
}

function closeTS3ChannelPicker() {
    const modal = document.getElementById('ts3-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

async function saveTS3Settings() {
    const payload = {
        host: document.getElementById('ts3-host').value.trim(),
        query_port: Number(document.getElementById('ts3-query-port').value),
        voice_port: Number(document.getElementById('ts3-voice-port').value),
        query_username: document.getElementById('ts3-query-user').value.trim(),
        query_password: document.getElementById('ts3-query-pass').value,
        bot_nickname: document.getElementById('ts3-bot-nickname').value.trim(),
        default_channel: document.getElementById('ts3-default-channel').value.trim(),
        query_slowmode_ms: Number(document.getElementById('ts3-query-slowmode').value)
    };

    const res = await ts3Fetch('/api/settings/ts3', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    });

    if (!res.ok) {
        setTS3Status(`${t('common.error')}: ${await res.text()}`, true);
        return;
    }

    const out = await res.json();
    if (out && out.restart_required) {
        setTS3Status(t('ts3.savedRestart'), false);
    } else {
        setTS3Status(t('ts3.saved'), false);
    }
}

document.addEventListener('DOMContentLoaded', loadTS3Settings);
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('ts3-channel-modal');
    if (!modal) {
        return;
    }

    modal.addEventListener('click', (ev) => {
        if (ev.target === modal) {
            closeTS3ChannelPicker();
        }
    });

    document.addEventListener('keydown', (ev) => {
        if (ev.key === 'Escape') {
            closeTS3ChannelPicker();
        }
    });
});
