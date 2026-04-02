const supportFetch = typeof authFetch === 'function' ? authFetch : fetch;

let supportChannels = [];          // all channels loaded from TS3
let supportSelectedChannels = [];  // array of {id: number, name: string}
let supportServerGroups = [];      // all server groups loaded from API
let supportSelectedGroups = [];    // array of {id: number, name: string}

function setSupportStatus(id, message, isError) {
    const el = document.getElementById(id);
    if (!el) {
        return;
    }
    el.textContent = message;
    el.className = 'save-status ' + (isError ? 'err' : 'ok');
}

// ---- Channel tag list ----

function renderSupportTags() {
    const wrap = document.getElementById('support-channel-tags');
    if (!wrap) {
        return;
    }
    if (!supportSelectedChannels.length) {
        wrap.innerHTML = '<span class="muted" style="font-size:0.9em">' + t('support.noChannelsSelected') + '</span>';
        return;
    }
    wrap.innerHTML = supportSelectedChannels.map((ch, i) => {
        const label = (ch.name && ch.name !== '#' + ch.id ? ch.name + ' ' : '') + '#' + ch.id;
        return '<span class="tag">' + label + '<button onclick="removeSupportChannel(' + i + ')" title="' + t('common.remove') + '">&times;</button></span>';
    }).join('');
}

function removeSupportChannel(index) {
    supportSelectedChannels.splice(index, 1);
    renderSupportTags();
    if (document.getElementById('support-channel-modal')?.classList.contains('open')) {
        renderSupportChannelPicker();
    }
}

function addSupportChannel(selectedId) {
    const input = document.getElementById('support-new-channel');
    const id = Number.isInteger(selectedId) ? selectedId : parseInt(input?.value || '0', 10);
    if (!Number.isInteger(id) || id <= 0) {
        return;
    }

    const existingIndex = supportSelectedChannels.findIndex((ch) => ch.id === id);
    if (existingIndex !== -1) {
        if (input) {
            input.value = '';
        }
        return;
    }

    const found = supportChannels.find((ch) => Number(ch.id) === id);
    supportSelectedChannels.push({
        id,
        name: found ? String(found.name || '').trim() : ''
    });

    renderSupportTags();
    if (document.getElementById('support-channel-modal')?.classList.contains('open')) {
        renderSupportChannelPicker();
    }

    if (input) {
        input.value = '';
    }
}

// ---- Channel picker modal ----

function renderSupportChannelPicker() {
    const list = document.getElementById('support-channel-picker-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(supportChannels) || supportChannels.length === 0) {
        list.innerHTML = '<div class="channel-empty">' + t('support.noChannels') + '</div>';
        return;
    }

    const selectedIds = new Set(supportSelectedChannels.map((ch) => ch.id));

    list.innerHTML = supportChannels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('support.unnamedChannel');
        const sel = selectedIds.has(id) ? ' selected' : '';
        return '<button type="button" class="channel-item' + sel + '" data-channel-id="' + id + '" data-channel-name="' + name.replace(/"/g, '&quot;') + '">' + name + '<small>#' + id + '</small></button>';
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const id = parseInt(btn.getAttribute('data-channel-id') || '0', 10);
            if (id <= 0) {
                return;
            }
            const idx = supportSelectedChannels.findIndex((ch) => ch.id === id);
            if (idx === -1) {
                addSupportChannel(id);
            } else {
                supportSelectedChannels.splice(idx, 1);
                renderSupportTags();
            }
            renderSupportChannelPicker();
        });
    });
}

async function openSupportChannelPicker() {
    const modal = document.getElementById('support-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (supportChannels.length > 0) {
        renderSupportChannelPicker();
        return;
    }

    const list = document.getElementById('support-channel-picker-list');
    if (list) {
        list.innerHTML = '<div class="channel-empty">' + t('common.loading') + '</div>';
    }

    try {
        const res = await supportFetch('/api/settings/ts3/channels');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('support.channelsLoadError'));
        }
        const payload = await res.json();
        supportChannels = Array.isArray(payload.channels) ? payload.channels : [];
        renderSupportChannelPicker();
    } catch (e) {
        if (list) {
            list.innerHTML = '<div class="channel-empty">' + (e && e.message ? e.message : t('support.channelsLoadError')) + '</div>';
        }
    }
}

function closeSupportChannelPicker() {
    const modal = document.getElementById('support-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function selectedSupportChannels() {
    return supportSelectedChannels.map((ch) => ch.id).filter((id) => Number.isInteger(id) && id > 0);
}

// ---- Supporter group tag list ----

function renderGroupTags() {
    const wrap = document.getElementById('support-supporter-group-tags');
    if (!wrap) {
        return;
    }
    if (!supportSelectedGroups.length) {
        wrap.innerHTML = '<span class="muted" style="font-size:0.9em">' + t('support.noGroupsSelected') + '</span>';
        return;
    }
    wrap.innerHTML = supportSelectedGroups.map((g, i) => {
        const label = (g.name && g.name !== '#' + g.id ? g.name + ' ' : '') + '#' + g.id;
        return '<span class="tag">' + label + '<button onclick="removeSupporterGroup(' + i + ')" title="' + t('common.remove') + '">&times;</button></span>';
    }).join('');
}

function removeSupporterGroup(index) {
    supportSelectedGroups.splice(index, 1);
    renderGroupTags();
    if (document.getElementById('support-group-modal')?.classList.contains('open')) {
        renderSupporterGroupPicker();
    }
}

function addSupporterGroup(selectedId) {
    const input = document.getElementById('support-new-group');
    const id = Number.isInteger(selectedId) ? selectedId : parseInt(input?.value || '0', 10);
    if (!Number.isInteger(id) || id <= 0) {
        return;
    }

    const existingIndex = supportSelectedGroups.findIndex((g) => g.id === id);
    if (existingIndex !== -1) {
        if (input) {
            input.value = '';
        }
        return;
    }

    const found = supportServerGroups.find((g) => Number(g.id) === id);
    supportSelectedGroups.push({
        id,
        name: found ? String(found.name || '').trim() : ''
    });

    renderGroupTags();
    if (document.getElementById('support-group-modal')?.classList.contains('open')) {
        renderSupporterGroupPicker();
    }

    if (input) {
        input.value = '';
    }
}

// ---- Supporter group picker modal ----

function renderSupporterGroupPicker() {
    const list = document.getElementById('support-group-picker-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(supportServerGroups) || supportServerGroups.length === 0) {
        list.innerHTML = '<div class="channel-empty">' + t('support.noServerGroups') + '</div>';
        return;
    }

    const selectedIds = new Set(supportSelectedGroups.map((g) => g.id));

    list.innerHTML = supportServerGroups.map((g) => {
        const id = Number(g.id) || 0;
        const name = String(g.name || '').trim() || t('support.unnamedGroup');
        const sel = selectedIds.has(id) ? ' selected' : '';
        return '<button type="button" class="channel-item' + sel + '" data-group-id="' + id + '" data-group-name="' + name.replace(/"/g, '&quot;') + '">' + name + '<small>#' + id + '</small></button>';
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const id = parseInt(btn.getAttribute('data-group-id') || '0', 10);
            if (id <= 0) {
                return;
            }
            const idx = supportSelectedGroups.findIndex((g) => g.id === id);
            if (idx === -1) {
                addSupporterGroup(id);
            } else {
                supportSelectedGroups.splice(idx, 1);
                renderGroupTags();
            }
            renderSupporterGroupPicker();
        });
    });
}

async function openSupporterGroupPicker() {
    const modal = document.getElementById('support-group-modal');
    if (!modal) {
        return;
    }
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (supportServerGroups.length > 0) {
        renderSupporterGroupPicker();
        return;
    }

    const list = document.getElementById('support-group-picker-list');
    if (list) {
        list.innerHTML = '<div class="channel-empty">' + t('common.loading') + '</div>';
    }

    try {
        await loadSupportServerGroups();
        renderSupporterGroupPicker();
    } catch (e) {
        if (list) {
            list.innerHTML = '<div class="channel-empty">' + (e && e.message ? e.message : t('support.serverGroupsLoadError')) + '</div>';
        }
    }
}

function closeSupporterGroupPicker() {
    const modal = document.getElementById('support-group-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function selectedSupporterGroups() {
    return supportSelectedGroups.map((g) => g.id).filter((id) => Number.isInteger(id) && id > 0);
}

function updateSupportWaitingPreview() {
    const preview = document.getElementById('support-waiting-preview');
    const input = document.getElementById('support-waiting-area');
    if (!preview || !input) {
        return;
    }

    const id = parseInt(input.value, 10);
    if (isNaN(id) || id <= 0) {
        preview.textContent = t('support.noWaitingArea');
        return;
    }

    const found = supportChannels.find((ch) => Number(ch.id) === id);
    if (found) {
        preview.textContent = (found.name || t('support.unnamedChannel')) + ' (#' + id + ')';
    } else {
        preview.textContent = '#' + id;
    }
}

function renderSupportWaitingAreaList() {
    const list = document.getElementById('support-waiting-list');
    if (!list) {
        return;
    }

    if (!Array.isArray(supportChannels) || supportChannels.length === 0) {
        list.innerHTML = '<div class="channel-empty">' + t('support.noChannels') + '</div>';
        return;
    }

    const waitingId = Number(document.getElementById('support-waiting-area')?.value || 0);

    list.innerHTML = supportChannels.map((ch) => {
        const id = Number(ch.id) || 0;
        const name = String(ch.name || '').trim() || t('support.unnamedChannel');
        const sel = waitingId === id ? ' selected' : '';
        return '<button type="button" class="channel-item' + sel + '" data-channel-id="' + id + '">' + name + '<small>#' + id + '</small></button>';
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const id = parseInt(btn.getAttribute('data-channel-id') || '0', 10);
            const input = document.getElementById('support-waiting-area');
            if (input) {
                input.value = String(id > 0 ? id : 0);
            }
            updateSupportWaitingPreview();
            closeSupportWaitingAreaPicker();
        });
    });
}

async function openSupportWaitingAreaPicker() {
    const modal = document.getElementById('support-waiting-modal');
    if (!modal) {
        return;
    }
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    if (!supportChannels.length) {
        try {
            await loadSupportChannels();
        } catch (e) {
            const list = document.getElementById('support-waiting-list');
            if (list) {
                list.innerHTML = '<div class="channel-empty">' + (e && e.message ? e.message : t('support.channelsLoadError')) + '</div>';
            }
            return;
        }
    }

    renderSupportWaitingAreaList();
}

function closeSupportWaitingAreaPicker() {
    const modal = document.getElementById('support-waiting-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function applySupportSettings(settings) {
    const ids = settings.support_channel_ids || [];
    supportSelectedChannels = ids.map((id) => {
        const ch = supportChannels.find((c) => Number(c.id) === id);
        return { id, name: ch ? (ch.name || t('support.unnamedChannel')) : '' };
    });
    renderSupportTags();

    const waitingAreaInput = document.getElementById('support-waiting-area');
    if (waitingAreaInput) {
        waitingAreaInput.value = String(settings.waiting_area_channel_id || 0);
    }
    updateSupportWaitingPreview();

    const openMessage = document.getElementById('support-open-message');
    const closeMessage = document.getElementById('support-close-message');
    const joinOpenMessage = document.getElementById('support-join-open-message');
    const joinCloseMessage = document.getElementById('support-join-close-message');
    const supporterMessage = document.getElementById('support-supporter-message');
    const autoEnabled = document.getElementById('support-auto-enabled');
    const openTime = document.getElementById('support-open-time');
    const closeTime = document.getElementById('support-close-time');

    if (openMessage) {
        openMessage.value = settings.open_poke_message || '';
    }
    if (closeMessage) {
        closeMessage.value = settings.closed_poke_message || '';
    }
    if (joinOpenMessage) {
        joinOpenMessage.value = settings.join_open_poke_message || '';
    }
    if (joinCloseMessage) {
        joinCloseMessage.value = settings.join_closed_poke_message || '';
    }
    if (supporterMessage) {
        supporterMessage.value = settings.supporter_poke_message || '';
    }
    const groupIds = settings.supporter_group_ids || [];
    supportSelectedGroups = groupIds.map((id) => {
        const g = supportServerGroups.find((sg) => Number(sg.id) === id);
        return { id, name: g ? (g.name || t('support.unnamedGroup')) : '' };
    });
    renderGroupTags();
    if (autoEnabled) {
        autoEnabled.checked = !!settings.auto_schedule_enabled;
    }
    if (openTime) {
        openTime.value = settings.auto_open_time || '08:00';
    }
    if (closeTime) {
        closeTime.value = settings.auto_close_time || '22:00';
    }
}

async function loadSupportChannels() {
    const res = await supportFetch('/api/settings/ts3/channels');
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || t('support.channelsLoadError'));
    }
    const payload = await res.json();
    supportChannels = Array.isArray(payload.channels) ? payload.channels : [];
}

async function loadSupportServerGroups() {
    const res = await supportFetch('/api/settings/ts3/servergroups');
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || t('support.serverGroupsLoadError'));
    }
    const payload = await res.json();
    supportServerGroups = Array.isArray(payload.groups) ? payload.groups : [];
}

async function loadSupportSettings() {
    const res = await supportFetch('/api/support/settings');
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || t('support.settingsLoadError'));
    }
    return res.json();
}

async function loadSupportStatus() {
    try {
        const res = await supportFetch('/api/support/status');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('support.statusLoadError'));
        }
        const payload = await res.json();
        const status = payload && payload.status ? payload.status : {};

        const state = document.getElementById('support-current-state');
        const action = document.getElementById('support-last-action');
        const error = document.getElementById('support-last-error');
        if (state) {
            state.textContent = status.open ? t('support.stateOpen') : t('support.stateClosed');
            state.style.color = status.open ? 'var(--success)' : 'var(--danger)';
        }
        if (action) {
            action.textContent = status.last_action || '-';
        }
        if (error) {
            error.textContent = status.last_error || '-';
        }
    } catch (err) {
        setSupportStatus('support-action-status', t('support.statusLoadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function saveSupportSettings() {
    const payload = {
        support_channel_ids: selectedSupportChannels(),
        waiting_area_channel_id: Number(document.getElementById('support-waiting-area')?.value || 0),
        open_poke_message: (document.getElementById('support-open-message')?.value || '').trim(),
        closed_poke_message: (document.getElementById('support-close-message')?.value || '').trim(),
        join_open_poke_message: (document.getElementById('support-join-open-message')?.value || '').trim(),
        join_closed_poke_message: (document.getElementById('support-join-close-message')?.value || '').trim(),
        supporter_poke_message: (document.getElementById('support-supporter-message')?.value || '').trim(),
        supporter_group_ids: selectedSupporterGroups(),
        auto_schedule_enabled: !!document.getElementById('support-auto-enabled')?.checked,
        auto_open_time: document.getElementById('support-open-time')?.value || '',
        auto_close_time: document.getElementById('support-close-time')?.value || ''
    };

    if (!payload.support_channel_ids.length) {
        setSupportStatus('support-settings-status', t('support.validationChannels'), true);
        return;
    }

    try {
        const res = await supportFetch('/api/support/settings', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('support.saveError'));
        }
        setSupportStatus('support-settings-status', t('support.saved'), false);
    } catch (err) {
        setSupportStatus('support-settings-status', t('support.saveError', { error: err && err.message ? err.message : err }), true);
    }
}

async function executeSupportAction(action) {
    try {
        const res = await supportFetch('/api/support/action', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ action })
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('support.actionError'));
        }
        const payload = await res.json();
        setSupportStatus('support-action-status', t('support.actionDone', { action }), false);

        const status = payload && payload.status ? payload.status : null;
        if (status) {
            const state = document.getElementById('support-current-state');
            const lastAction = document.getElementById('support-last-action');
            const lastError = document.getElementById('support-last-error');
            if (state) {
                state.textContent = status.open ? t('support.stateOpen') : t('support.stateClosed');
                state.style.color = status.open ? 'var(--success)' : 'var(--danger)';
            }
            if (lastAction) {
                lastAction.textContent = status.last_action || '-';
            }
            if (lastError) {
                lastError.textContent = status.last_error || '-';
            }
        }
    } catch (err) {
        setSupportStatus('support-action-status', t('support.actionError', { error: err && err.message ? err.message : err }), true);
    }
}

async function initSupportPage() {
    try {
        // Load plugin active state for badge
        const pluginsRes = await supportFetch('/api/plugins');
        if (pluginsRes.ok) {
            const plugins = await pluginsRes.json();
            const sc = plugins.find((p) => p.name === 'SupportControl');
            const badge = document.getElementById('support-status-badge');
            if (sc && badge) {
                badge.textContent = sc.active ? t('common.active') : t('common.inactive');
                badge.className = 'plugin-status-badge ' + (sc.active ? 'loaded' : 'unloaded');
            }
        }
    } catch (_) {
        // badge stays empty if plugins API fails
    }

    try {
        await loadSupportChannels();
        await loadSupportServerGroups();
        const settings = await loadSupportSettings();
        applySupportSettings(settings);
        setSupportStatus('support-settings-status', t('support.loaded'), false);
    } catch (err) {
        setSupportStatus('support-settings-status', t('support.settingsLoadError', { error: err && err.message ? err.message : err }), true);
    }

    await loadSupportStatus();

    const waitingAreaInput = document.getElementById('support-waiting-area');
    if (waitingAreaInput) {
        waitingAreaInput.addEventListener('input', updateSupportWaitingPreview);
    }

    const waitingModal = document.getElementById('support-waiting-modal');
    if (waitingModal) {
        waitingModal.addEventListener('click', (ev) => {
            if (ev.target === waitingModal) {
                closeSupportWaitingAreaPicker();
            }
        });
    }

    const groupModal = document.getElementById('support-group-modal');
    if (groupModal) {
        groupModal.addEventListener('click', (ev) => {
            if (ev.target === groupModal) {
                closeSupporterGroupPicker();
            }
        });
    }

    document.addEventListener('keydown', (ev) => {
        if (ev.key === 'Escape') {
            closeSupportWaitingAreaPicker();
            closeSupporterGroupPicker();
        }
    });
}

document.addEventListener('DOMContentLoaded', initSupportPage);
