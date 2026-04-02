let excludedChannels = [];
const secureFetch = typeof authFetch === 'function' ? authFetch : fetch;
let afkPluginActive = null;
let afkChannelsCache = [];
let afkTargetChannelName = '';

function escapeHtml(value) {
    return String(value ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function escapeAttr(value) {
    return escapeHtml(value).replace(/`/g, '&#96;');
}

async function loadAfkConfig() {
    try {
        // Plugin-Status laden
        const pluginsRes = await secureFetch('/api/plugins');
        if (pluginsRes.ok) {
            const plugins = await pluginsRes.json();
            const afk = plugins.find(p => p.name === 'AfkMover');
            const badge = document.getElementById('afk-status-badge');
            if (afk && badge) {
                afkPluginActive = !!afk.active;
                badge.textContent = afk.active ? t('common.active') : t('common.inactive');
                badge.className = 'plugin-status-badge ' + (afk.active ? 'loaded' : 'unloaded');
            }
        }

        // Konfiguration laden
        const res = await secureFetch('/api/plugins/config?name=AfkMover');
        if (!res.ok) {
            document.getElementById('afk-channel-id').value = 0;
            document.getElementById('afk-timeout').value = 10;
            document.getElementById('afk-return-on-activity').checked = false;
            excludedChannels = [];
            renderExcludedChannels();
            showStatus(t('afk.pluginDisabledInfo'), false);
            return;
        }
        const cfg = await res.json();
        excludedChannels = cfg.excluded_channels || [];
        afkTargetChannelName = String(cfg.afk_channel_name || '').trim();
        document.getElementById('afk-channel-id').value = cfg.afk_channel_id || 0;
        document.getElementById('afk-timeout').value = Math.round((cfg.timeout_seconds || 600) / 60);
        document.getElementById('afk-return-on-activity').checked = cfg.return_on_activity === true;
        updateAfkTargetPreview();
        renderExcludedChannels();
    } catch (e) {
        showStatus(t('afk.loadError', { error: e && e.message ? e.message : e }), true);
    }
}

function updateAfkTargetPreview() {
    const preview = document.getElementById('afk-target-preview');
    const input = document.getElementById('afk-channel-id');
    if (!preview || !input) {
        return;
    }

    const id = parseInt(input.value, 10);
    if (isNaN(id) || id <= 0) {
        preview.textContent = '-';
        return;
    }

    if (afkTargetChannelName) {
        preview.textContent = `${afkTargetChannelName} (#${id})`;
        return;
    }

    const fromCache = afkChannelsCache.find((ch) => Number(ch.id) === id);
    if (fromCache) {
        preview.textContent = `${fromCache.name || t('afk.unnamedChannel')} (#${id})`;
    } else {
        preview.textContent = `#${id}`;
    }
}

function renderAfkChannelList(channels) {
    const list = document.getElementById('afk-channel-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(channels) || channels.length === 0) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('afk.noChannels'))}</div>`;
        return;
    }

    list.innerHTML = channels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('afk.unnamedChannel');
        return `<button type="button" class="channel-item" data-channel-id="${id}" data-channel-name="${escapeAttr(name)}">${escapeHtml(name)}<small>#${id}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const selectedId = parseInt(btn.getAttribute('data-channel-id') || '0', 10);
            const selectedName = btn.getAttribute('data-channel-name') || '';
            if (selectedId > 0) {
                document.getElementById('afk-channel-id').value = selectedId;
                afkTargetChannelName = selectedName;
                updateAfkTargetPreview();
            }
            closeAfkChannelPicker();
        });
    });
}

async function openAfkChannelPicker() {
    const modal = document.getElementById('afk-channel-modal');
    if (!modal) {
        return;
    }

    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (afkChannelsCache.length > 0) {
        renderAfkChannelList(afkChannelsCache);
        return;
    }

    const list = document.getElementById('afk-channel-list');
    if (list) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('common.loading'))}</div>`;
    }

    try {
        const res = await secureFetch('/api/settings/ts3/channels');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('afk.channelLoadError'));
        }
        const payload = await res.json();
        afkChannelsCache = Array.isArray(payload.channels) ? payload.channels : [];
        renderAfkChannelList(afkChannelsCache);
    } catch (e) {
        if (list) {
            const msg = e && e.message ? e.message : t('afk.channelLoadError');
            list.innerHTML = `<div class="channel-empty">${escapeHtml(msg)}</div>`;
        }
    }
}

function closeAfkChannelPicker() {
    const modal = document.getElementById('afk-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function renderExcludedChannelList(channels) {
    const list = document.getElementById('excluded-channel-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(channels) || channels.length === 0) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('afk.noChannels'))}</div>`;
        return;
    }

    list.innerHTML = channels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('afk.unnamedChannel');
        return `<button type="button" class="channel-item" data-channel-id="${id}">${escapeHtml(name)}<small>#${id}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const selectedId = parseInt(btn.getAttribute('data-channel-id') || '0', 10);
            if (selectedId > 0) {
                addExcludedChannel(selectedId);
            }
            closeExcludedChannelPicker();
        });
    });
}

async function openExcludedChannelPicker() {
    const modal = document.getElementById('excluded-channel-modal');
    if (!modal) {
        return;
    }

    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (afkChannelsCache.length > 0) {
        renderExcludedChannelList(afkChannelsCache);
        return;
    }

    const list = document.getElementById('excluded-channel-list');
    if (list) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('common.loading'))}</div>`;
    }

    try {
        const res = await secureFetch('/api/settings/ts3/channels');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('afk.channelLoadError'));
        }
        const payload = await res.json();
        afkChannelsCache = Array.isArray(payload.channels) ? payload.channels : [];
        renderExcludedChannelList(afkChannelsCache);
    } catch (e) {
        if (list) {
            const msg = e && e.message ? e.message : t('afk.channelLoadError');
            list.innerHTML = `<div class="channel-empty">${escapeHtml(msg)}</div>`;
        }
    }
}

function closeExcludedChannelPicker() {
    const modal = document.getElementById('excluded-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function renderExcludedChannels() {
    const list = document.getElementById('excluded-channels-list');
    if (!excludedChannels.length) {
        list.innerHTML = `<span class="muted" style="font-size:0.9em">${escapeHtml(t('afk.noExceptions'))}</span>`;
        return;
    }
    list.innerHTML = excludedChannels.map((ch, i) =>
        `<span class="tag">${escapeHtml(ch)}<button onclick="removeExcludedChannel(${i})" title="${escapeAttr(t('common.remove'))}">&times;</button></span>`
    ).join('');
}

function addExcludedChannel(selectedId) {
    const input = document.getElementById('new-excluded-channel');
    const val = Number.isInteger(selectedId) ? selectedId : parseInt(input.value, 10);
    if (!isNaN(val) && val > 0 && !excludedChannels.includes(val)) {
        excludedChannels.push(val);
        renderExcludedChannels();
        if (input) {
            input.value = '';
        }
    }
}

function removeExcludedChannel(i) {
    excludedChannels.splice(i, 1);
    renderExcludedChannels();
}

async function saveAfkConfig() {
    const enabled = afkPluginActive === true;

    const timeoutMin = parseInt(document.getElementById('afk-timeout').value) || 10;
    const config = {
        enabled: enabled,
        afk_channel_id: parseInt(document.getElementById('afk-channel-id').value) || 0,
        afk_channel_name: afkTargetChannelName,
        timeout_seconds: timeoutMin * 60,
        return_on_activity: document.getElementById('afk-return-on-activity').checked === true,
        excluded_channels: excludedChannels
    };
    const res = await secureFetch('/api/plugins/config?name=AfkMover', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
    });
    if (res.ok) {
        showStatus(t('common.saved'), false);
    } else {
        const err = await res.text();
        showStatus(`${t('common.error')}: ${err}`, true);
    }
}

function showStatus(msg, isErr) {
    const el = document.getElementById('save-status');
    el.textContent = msg;
    el.className = 'save-status ' + (isErr ? 'err' : 'ok');
    setTimeout(() => { el.textContent = ''; el.className = 'save-status'; }, 4000);
}

document.addEventListener('DOMContentLoaded', loadAfkConfig);
document.addEventListener('DOMContentLoaded', () => {
    const input = document.getElementById('afk-channel-id');
    if (input) {
        input.addEventListener('input', () => {
            afkTargetChannelName = '';
            updateAfkTargetPreview();
        });
    }

    const modal = document.getElementById('afk-channel-modal');
    if (modal) {
        modal.addEventListener('click', (ev) => {
            if (ev.target === modal) {
                closeAfkChannelPicker();
            }
        });
    }

    const excludedModal = document.getElementById('excluded-channel-modal');
    if (excludedModal) {
        excludedModal.addEventListener('click', (ev) => {
            if (ev.target === excludedModal) {
                closeExcludedChannelPicker();
            }
        });
    }

    document.addEventListener('keydown', (ev) => {
        if (ev.key === 'Escape') {
            closeAfkChannelPicker();
            closeExcludedChannelPicker();
        }
    });
});
