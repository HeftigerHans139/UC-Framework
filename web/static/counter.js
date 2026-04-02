let adminGroups = [];
let excludedGroups = [];
let excludedNicks = [];
let counterChannelsCache = [];
let activeCounterPickerTarget = null;
const secureFetch = typeof authFetch === 'function' ? authFetch : fetch;

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

async function loadCounterConfigs() {
    try {
        // Plugin-Status laden
        const pluginsRes = await secureFetch('/api/plugins');
        if (pluginsRes.ok) {
            const plugins = await pluginsRes.json();
            updateBadge('admin-status-badge', plugins.find(p => p.name === 'AdminCounter'));
            updateBadge('member-status-badge', plugins.find(p => p.name === 'MemberCounter'));
        }

        // Admin Counter Konfiguration
        const adminRes = await secureFetch('/api/plugins/config?name=AdminCounter');
        if (adminRes.ok) {
            const cfg = await adminRes.json();
            adminGroups = cfg.admin_groups || [];
            document.getElementById('admin-rename-channel-id').value = cfg.rename_channel_id || 0;
            document.getElementById('admin-rename-template').value = cfg.rename_name_template || '';
            document.getElementById('admin-rename-token').value = cfg.rename_count_token || '{count}';
            renderAdminGroups();
        } else {
            document.getElementById('admin-save-status').textContent = t('counter.pluginNotLoaded');
            document.getElementById('admin-save-status').className = 'save-status err';
        }

        // Member Counter Konfiguration
        const memberRes = await secureFetch('/api/plugins/config?name=MemberCounter');
        if (memberRes.ok) {
            const cfg = await memberRes.json();
            excludedGroups = cfg.excluded_groups || [];
            excludedNicks = cfg.excluded_nicknames || [];
            document.getElementById('member-rename-channel-id').value = cfg.rename_channel_id || 0;
            document.getElementById('member-rename-template').value = cfg.rename_name_template || '';
            document.getElementById('member-rename-token').value = cfg.rename_count_token || '{count}';
            renderExcludedGroups();
            renderExcludedNicks();
        } else {
            document.getElementById('member-save-status').textContent = t('counter.pluginNotLoaded');
            document.getElementById('member-save-status').className = 'save-status err';
        }
    } catch (e) {
        console.error('Fehler beim Laden:', e);
    }
}

function updateBadge(id, plugin) {
    const badge = document.getElementById(id);
    if (badge && plugin) {
        badge.textContent = plugin.active ? t('common.active') : t('common.inactive');
        badge.className = 'plugin-status-badge ' + (plugin.active ? 'loaded' : 'unloaded');
    }
}

function renderTagList(containerId, items, removeFnName) {
    const el = document.getElementById(containerId);
    if (!items || !items.length) {
        el.innerHTML = `<span style="color:#aaa;font-size:0.9em">${escapeHtml(t('common.noEntries'))}</span>`;
        return;
    }
    el.innerHTML = items.map((item, i) =>
        `<span class="tag">${escapeHtml(item)}<button onclick="${removeFnName}(${i})" title="${escapeAttr(t('common.remove'))}">&times;</button></span>`
    ).join('');
}

function renderAdminGroups()    { renderTagList('admin-groups-list',    adminGroups,    'removeAdminGroup'); }
function renderExcludedGroups() { renderTagList('excluded-groups-list', excludedGroups, 'removeExcludedGroup'); }
function renderExcludedNicks()  { renderTagList('excluded-nicks-list',  excludedNicks,  'removeExcludedNick'); }

function addAdminGroup() {
    const val = parseInt(document.getElementById('new-admin-group').value);
    if (!isNaN(val) && val > 0 && !adminGroups.includes(val)) {
        adminGroups.push(val);
        renderAdminGroups();
        document.getElementById('new-admin-group').value = '';
    }
}
function removeAdminGroup(i) { adminGroups.splice(i, 1); renderAdminGroups(); }

function addExcludedGroup() {
    const val = parseInt(document.getElementById('new-excluded-group').value);
    if (!isNaN(val) && val > 0 && !excludedGroups.includes(val)) {
        excludedGroups.push(val);
        renderExcludedGroups();
        document.getElementById('new-excluded-group').value = '';
    }
}
function removeExcludedGroup(i) { excludedGroups.splice(i, 1); renderExcludedGroups(); }

function addExcludedNick() {
    const val = document.getElementById('new-excluded-nick').value.trim();
    if (val && !excludedNicks.includes(val)) {
        excludedNicks.push(val);
        renderExcludedNicks();
        document.getElementById('new-excluded-nick').value = '';
    }
}
function removeExcludedNick(i) { excludedNicks.splice(i, 1); renderExcludedNicks(); }

function setCounterChannelPickerVisible(isVisible) {
    ['admin-rename-channel-picker', 'member-rename-channel-picker'].forEach((id) => {
        const btn = document.getElementById(id);
        if (!btn) {
            return;
        }
        btn.style.display = isVisible ? 'inline-flex' : 'none';
    });
}

function renderCounterChannelList(channels) {
    const list = document.getElementById('counter-channel-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(channels) || channels.length === 0) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('ts3.noChannels'))}</div>`;
        return;
    }

    list.innerHTML = channels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('ts3.unnamedChannel');
        return `<button type="button" class="channel-item" data-channel-id="${id}" data-target="${escapeAttr(activeCounterPickerTarget || '')}">${escapeHtml(name)}<small>#${id}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', async () => {
            const selectedId = parseInt(btn.getAttribute('data-channel-id') || '0', 10) || 0;
            const selectedTarget = btn.getAttribute('data-target') || '';
            const targetInputId = selectedTarget === 'member' ? 'member-rename-channel-id' : 'admin-rename-channel-id';
            const targetInput = document.getElementById(targetInputId);
            if (targetInput) {
                targetInput.value = selectedId;
            }

            closeCounterChannelPicker();

            if (selectedTarget === 'member') {
                await saveMemberConfig();
            } else {
                await saveAdminConfig();
            }
        });
    });
}

async function loadCounterChannelPickerAvailability() {
    try {
        const res = await secureFetch('/api/ts3/connection/status');
        if (!res.ok) {
            setCounterChannelPickerVisible(false);
            return;
        }

        const payload = await res.json();
        const status = payload && payload.status ? payload.status : {};
        const isConnected = status.bot_running !== false && !!status.connected;
        setCounterChannelPickerVisible(isConnected);
    } catch (_) {
        setCounterChannelPickerVisible(false);
    }
}

async function openCounterChannelPicker(target) {
    const modal = document.getElementById('counter-channel-modal');
    const list = document.getElementById('counter-channel-list');
    if (!modal || !list) {
        return;
    }

    activeCounterPickerTarget = target === 'member' ? 'member' : 'admin';
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (counterChannelsCache.length > 0) {
        renderCounterChannelList(counterChannelsCache);
        return;
    }

    list.innerHTML = `<div class="channel-empty">${escapeHtml(t('common.loading'))}</div>`;
    try {
        const res = await secureFetch('/api/settings/ts3/channels');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('ts3.channelLoadError'));
        }

        const payload = await res.json();
        counterChannelsCache = Array.isArray(payload.channels) ? payload.channels : [];
        renderCounterChannelList(counterChannelsCache);
    } catch (err) {
        const msg = err && err.message ? err.message : t('ts3.channelLoadError');
        list.innerHTML = `<div class="channel-empty">${escapeHtml(msg)}</div>`;
    }
}

function closeCounterChannelPicker() {
    const modal = document.getElementById('counter-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

async function postConfig(name, body, statusId) {
    const status = document.getElementById(statusId);
    const res = await secureFetch('/api/plugins/config?name=' + name, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    });
    if (res.ok) {
        let persisted = null;
        try {
            persisted = await res.json();
        } catch (_) {
            persisted = null;
        }

        if (persisted && Object.prototype.hasOwnProperty.call(body, 'rename_name_template')) {
            const expectedTemplate = body.rename_name_template || '';
            const expectedToken = body.rename_count_token || '{count}';
            const actualTemplate = persisted.rename_name_template || '';
            const actualToken = persisted.rename_count_token || '{count}';

            if (expectedTemplate !== actualTemplate || expectedToken !== actualToken) {
                status.textContent = `${t('common.error')}: ${t('counter.saveVerifyFailed')}`;
                status.className = 'save-status err';
                setTimeout(() => { status.textContent = ''; status.className = 'save-status'; }, 5000);
                return false;
            }
        }

        status.textContent = t('common.saved');
        status.className = 'save-status ok';
        setTimeout(() => { status.textContent = ''; status.className = 'save-status'; }, 4000);
        return true;
    } else {
        status.textContent = `${t('common.error')}: ${await res.text()}`;
        status.className = 'save-status err';
        setTimeout(() => { status.textContent = ''; status.className = 'save-status'; }, 4000);
        return false;
    }
}

function saveAdminConfig() {
    return postConfig('AdminCounter', {
        admin_groups: adminGroups,
        rename_channel_id: parseInt(document.getElementById('admin-rename-channel-id').value, 10) || 0,
        rename_name_template: document.getElementById('admin-rename-template').value.trim(),
        rename_count_token: document.getElementById('admin-rename-token').value.trim() || '{count}'
    }, 'admin-save-status');
}

function saveMemberConfig() {
    return postConfig('MemberCounter', {
        excluded_groups: excludedGroups,
        excluded_nicknames: excludedNicks,
        rename_channel_id: parseInt(document.getElementById('member-rename-channel-id').value, 10) || 0,
        rename_name_template: document.getElementById('member-rename-template').value.trim(),
        rename_count_token: document.getElementById('member-rename-token').value.trim() || '{count}'
    }, 'member-save-status');
}

document.addEventListener('DOMContentLoaded', loadCounterConfigs);
document.addEventListener('DOMContentLoaded', loadCounterChannelPickerAvailability);
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('counter-channel-modal');
    if (!modal) {
        return;
    }

    modal.addEventListener('click', (ev) => {
        if (ev.target === modal) {
            closeCounterChannelPicker();
        }
    });

    document.addEventListener('keydown', (ev) => {
        if (ev.key === 'Escape') {
            closeCounterChannelPicker();
        }
    });
});
