const discordFetch = typeof authFetch === 'function' ? authFetch : fetch;

const discordRoleState = {
    admin: [],
    supporter: [],
    bot: []
};

let discordChannelsCache = [];
let discordRolesCache = [];
let discordActiveChannelTarget = '';
let discordAllowedChannelTypes = [];
let discordActiveRoleKind = 'admin';

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

function setDiscordStatus(message, isError = false) {
    const el = document.getElementById('discord-settings-status');
    if (!el) {
        return;
    }
    el.textContent = message;
    el.className = 'save-status ' + (isError ? 'err' : 'ok');
}

function roleInputId(kind) {
    return `discord-${kind}-role-input`;
}

function roleWrapId(kind) {
    return `discord-${kind}-role-tags`;
}

function renderDiscordRoleTags(kind) {
    const wrap = document.getElementById(roleWrapId(kind));
    if (!wrap) {
        return;
    }
    const values = Array.isArray(discordRoleState[kind]) ? discordRoleState[kind] : [];
    if (!values.length) {
        wrap.innerHTML = `<span class="muted" style="font-size:0.9em">${escapeHtml(t('discord.noRolesSelected'))}</span>`;
        return;
    }
    wrap.innerHTML = values.map((value, index) => `<span class="tag">#${escapeHtml(value)}<button onclick="removeDiscordRole('${kind}', ${index})" title="${escapeAttr(t('common.remove'))}">&times;</button></span>`).join('');
}

function removeDiscordRole(kind, index) {
    discordRoleState[kind].splice(index, 1);
    renderDiscordRoleTags(kind);
}

function addDiscordRole(kind) {
    const input = document.getElementById(roleInputId(kind));
    const value = (input?.value || '').trim();
    if (!value) {
        return;
    }
    if (!discordRoleState[kind].includes(value)) {
        discordRoleState[kind].push(value);
    }
    renderDiscordRoleTags(kind);
    if (input) {
        input.value = '';
    }
}

function discordChannelLabel(channel) {
    const name = String(channel.name || '').trim() || t('discord.unnamedChannel');
    const type = String(channel.type || 'unknown').trim();
    return `${name} (${type})`;
}

function renderDiscordChannelPicker() {
    const list = document.getElementById('discord-channel-picker-list');
    if (!list) {
        return;
    }

    let channels = Array.isArray(discordChannelsCache) ? discordChannelsCache.slice() : [];
    if (discordAllowedChannelTypes.length) {
        const allowed = new Set(discordAllowedChannelTypes);
        channels = channels.filter((channel) => allowed.has(String(channel.type || '').trim()));
    }

    if (!channels.length) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('discord.noChannels'))}</div>`;
        return;
    }

    const targetValue = document.getElementById(discordActiveChannelTarget)?.value || '';
    list.innerHTML = channels.map((channel) => {
        const id = String(channel.id || '').trim();
        const name = discordChannelLabel(channel);
        const selected = targetValue === id ? ' selected' : '';
        return `<button type="button" class="channel-item${selected}" data-channel-id="${escapeAttr(id)}" data-channel-name="${escapeAttr(name)}">${escapeHtml(name)}<small>#${escapeHtml(id)}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-channel-id') || '';
            const input = document.getElementById(discordActiveChannelTarget);
            if (input) {
                input.value = id;
            }
            closeDiscordChannelPicker();
        });
    });
}

async function ensureDiscordChannelsLoaded() {
    if (discordChannelsCache.length > 0) {
        return;
    }
    const res = await discordFetch('/api/settings/discord/channels');
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || t('discord.channelsLoadError'));
    }
    const payload = await res.json();
    discordChannelsCache = Array.isArray(payload.channels) ? payload.channels : [];
}

async function openDiscordChannelPicker(targetId, allowedTypes) {
    discordActiveChannelTarget = targetId;
    discordAllowedChannelTypes = Array.isArray(allowedTypes) ? allowedTypes.slice() : [];

    const modal = document.getElementById('discord-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    const list = document.getElementById('discord-channel-picker-list');
    if (list) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('common.loading'))}</div>`;
    }

    try {
        await ensureDiscordChannelsLoaded();
        renderDiscordChannelPicker();
    } catch (err) {
        if (list) {
            list.innerHTML = `<div class="channel-empty">${escapeHtml(err && err.message ? err.message : t('discord.channelsLoadError'))}</div>`;
        }
    }
}

function closeDiscordChannelPicker() {
    const modal = document.getElementById('discord-channel-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function renderDiscordRolePicker() {
    const list = document.getElementById('discord-role-picker-list');
    if (!list) {
        return;
    }
    if (!discordRolesCache.length) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('discord.noRoles'))}</div>`;
        return;
    }

    const selected = new Set(discordRoleState[discordActiveRoleKind] || []);
    list.innerHTML = discordRolesCache.map((role) => {
        const id = String(role.id || '').trim();
        const name = String(role.name || '').trim() || t('discord.unnamedRole');
        const isSelected = selected.has(id) ? ' selected' : '';
        return `<button type="button" class="channel-item${isSelected}" data-role-id="${escapeAttr(id)}">${escapeHtml(name)}<small>#${escapeHtml(id)}</small></button>`;
    }).join('');

    list.querySelectorAll('.channel-item').forEach((btn) => {
        btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-role-id') || '';
            if (!id) {
                return;
            }
            const values = discordRoleState[discordActiveRoleKind] || [];
            const index = values.indexOf(id);
            if (index === -1) {
                values.push(id);
            } else {
                values.splice(index, 1);
            }
            discordRoleState[discordActiveRoleKind] = values;
            renderDiscordRoleTags(discordActiveRoleKind);
            renderDiscordRolePicker();
        });
    });
}

async function ensureDiscordRolesLoaded() {
    if (discordRolesCache.length > 0) {
        return;
    }
    const res = await discordFetch('/api/settings/discord/roles');
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || t('discord.rolesLoadError'));
    }
    const payload = await res.json();
    discordRolesCache = Array.isArray(payload.roles) ? payload.roles : [];
}

async function openDiscordRolePicker(kind) {
    discordActiveRoleKind = kind;
    const modal = document.getElementById('discord-role-modal');
    if (!modal) {
        return;
    }
    modal.classList.add('open');
    modal.setAttribute('aria-hidden', 'false');

    const list = document.getElementById('discord-role-picker-list');
    if (list) {
        list.innerHTML = `<div class="channel-empty">${escapeHtml(t('common.loading'))}</div>`;
    }

    try {
        await ensureDiscordRolesLoaded();
        renderDiscordRolePicker();
    } catch (err) {
        if (list) {
            list.innerHTML = `<div class="channel-empty">${escapeHtml(err && err.message ? err.message : t('discord.rolesLoadError'))}</div>`;
        }
    }
}

function closeDiscordRolePicker() {
    const modal = document.getElementById('discord-role-modal');
    if (!modal) {
        return;
    }
    modal.classList.remove('open');
    modal.setAttribute('aria-hidden', 'true');
}

function applyDiscordSettings(settings) {
    document.getElementById('discord-enabled').checked = !!settings.enabled;
    document.getElementById('discord-afk-kick-enabled').checked = !!settings.afk_kick_enabled;
    document.getElementById('discord-bot-token').value = settings.bot_token || '';
    document.getElementById('discord-application-id').value = settings.application_id || '';
    document.getElementById('discord-guild-id').value = settings.guild_id || '';
    document.getElementById('discord-afk-inactivity-minutes').value = settings.afk_inactivity_minutes > 0 ? settings.afk_inactivity_minutes : 30;
    document.getElementById('discord-bot-display-name').value = settings.bot_display_name || '';
    document.getElementById('discord-status-text').value = settings.status_text || '';
    document.getElementById('discord-command-prefix').value = settings.command_prefix || '!';
    document.getElementById('discord-log-channel-id').value = settings.log_channel_id || '';
    document.getElementById('discord-announcement-channel-id').value = settings.announcement_channel_id || '';
    document.getElementById('discord-support-category-id').value = settings.support_category_id || '';
    document.getElementById('discord-support-log-channel-id').value = settings.support_log_channel_id || '';

    discordRoleState.admin = Array.isArray(settings.admin_role_ids) ? settings.admin_role_ids.slice() : [];
    discordRoleState.supporter = Array.isArray(settings.supporter_role_ids) ? settings.supporter_role_ids.slice() : [];
    discordRoleState.bot = Array.isArray(settings.bot_role_ids) ? settings.bot_role_ids.slice() : [];

    renderDiscordRoleTags('admin');
    renderDiscordRoleTags('supporter');
    renderDiscordRoleTags('bot');
}

async function loadDiscordSettings() {
    try {
        const res = await discordFetch('/api/settings/discord');
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('discord.loadError'));
        }
        const settings = await res.json();
        applyDiscordSettings(settings);
        setDiscordStatus(t('discord.loaded'), false);
    } catch (err) {
        setDiscordStatus(t('discord.loadError', { error: err && err.message ? err.message : err }), true);
    }
}

async function saveDiscordSettings() {
    const payload = {
        enabled: !!document.getElementById('discord-enabled')?.checked,
        afk_kick_enabled: !!document.getElementById('discord-afk-kick-enabled')?.checked,
        bot_token: (document.getElementById('discord-bot-token')?.value || '').trim(),
        application_id: (document.getElementById('discord-application-id')?.value || '').trim(),
        guild_id: (document.getElementById('discord-guild-id')?.value || '').trim(),
        afk_inactivity_minutes: parseInt(document.getElementById('discord-afk-inactivity-minutes')?.value, 10) || 30,
        bot_display_name: (document.getElementById('discord-bot-display-name')?.value || '').trim(),
        status_text: (document.getElementById('discord-status-text')?.value || '').trim(),
        command_prefix: (document.getElementById('discord-command-prefix')?.value || '').trim(),
        log_channel_id: (document.getElementById('discord-log-channel-id')?.value || '').trim(),
        announcement_channel_id: (document.getElementById('discord-announcement-channel-id')?.value || '').trim(),
        support_category_id: (document.getElementById('discord-support-category-id')?.value || '').trim(),
        support_log_channel_id: (document.getElementById('discord-support-log-channel-id')?.value || '').trim(),
        admin_role_ids: discordRoleState.admin.slice(),
        supporter_role_ids: discordRoleState.supporter.slice(),
        bot_role_ids: discordRoleState.bot.slice()
    };

    try {
        const res = await discordFetch('/api/settings/discord', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || t('discord.saveError'));
        }
        const out = await res.json();
        if (out && out.restart_required) {
            setDiscordStatus(t('discord.savedRestart'), false);
            return;
        }
        setDiscordStatus(t('discord.saved'), false);
    } catch (err) {
        setDiscordStatus(t('discord.saveError', { error: err && err.message ? err.message : err }), true);
    }
}

document.addEventListener('DOMContentLoaded', loadDiscordSettings);
document.addEventListener('DOMContentLoaded', () => {
    const channelModal = document.getElementById('discord-channel-modal');
    if (channelModal) {
        channelModal.addEventListener('click', (ev) => {
            if (ev.target === channelModal) {
                closeDiscordChannelPicker();
            }
        });
    }

    const roleModal = document.getElementById('discord-role-modal');
    if (roleModal) {
        roleModal.addEventListener('click', (ev) => {
            if (ev.target === roleModal) {
                closeDiscordRolePicker();
            }
        });
    }

    document.addEventListener('keydown', (ev) => {
        if (ev.key === 'Escape') {
            closeDiscordChannelPicker();
            closeDiscordRolePicker();
        }
    });
});