// Dashboard-Logik für dynamische Karten und Plugin-Management
const secureFetch = typeof authFetch === 'function' ? authFetch : fetch;
let cachedPlugins = [];
let cachedTS3Settings = null;
let botStatusInterval = null;

function escapeHtml(value) {
    return String(value ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function renderPluginStatusList(plugins) {
    const container = document.getElementById('plugin-list');
    if (!container) {
        return;
    }

    container.innerHTML = '';

    plugins.forEach((plugin) => {
        const pluginLabel = typeof getPluginLabel === 'function' ? getPluginLabel(plugin.name) : plugin.name;
        const pluginDescription = typeof getPluginDescription === 'function' ? getPluginDescription(plugin) : (plugin.description || '');
        const row = document.createElement('div');
        row.className = 'plugin-row';
        row.innerHTML = `
            <div class="plugin-row-meta">
                <span class="plugin-row-name">${pluginLabel}</span>
                <span class="plugin-row-description">${pluginDescription}</span>
            </div>
            <span class="plugin-status-badge ${plugin.active ? 'loaded' : 'unloaded'}">${plugin.active ? t('common.active') : t('common.inactive')}</span>
        `;
        container.appendChild(row);
    });
}

function renderSettingsSummary(settings) {
    const container = document.getElementById('settings-content');
    if (!container) {
        return;
    }

    if (!settings || typeof settings !== 'object') {
        container.innerHTML = `<span class="plugin-row-description">${escapeHtml(t('dashboard.settingsUnavailable'))}</span>`;
        return;
    }

    const host = settings.host || '-';
    const queryPort = settings.query_port || 10011;
    const nickname = settings.bot_nickname || '-';
    const defaultChannel = settings.default_channel || '-';

    container.innerHTML = `
        <div class="plugin-row">
            <div class="plugin-row-meta">
                <span class="plugin-row-name">${escapeHtml(t('dashboard.ts3Host'))}</span>
                <span class="plugin-row-description">${escapeHtml(host)}</span>
            </div>
        </div>
        <div class="plugin-row">
            <div class="plugin-row-meta">
                <span class="plugin-row-name">${escapeHtml(t('dashboard.ts3QueryPort'))}</span>
                <span class="plugin-row-description">${escapeHtml(queryPort)}</span>
            </div>
        </div>
        <div class="plugin-row">
            <div class="plugin-row-meta">
                <span class="plugin-row-name">${escapeHtml(t('dashboard.ts3Nickname'))}</span>
                <span class="plugin-row-description">${escapeHtml(nickname)}</span>
            </div>
        </div>
        <div class="plugin-row">
            <div class="plugin-row-meta">
                <span class="plugin-row-name">${escapeHtml(t('dashboard.ts3DefaultChannel'))}</span>
                <span class="plugin-row-description">${escapeHtml(defaultChannel)}</span>
            </div>
        </div>
    `;
}

function renderBotStatusCards(bot, watchdog) {
    const container = document.getElementById('bot-status-cards');
    if (!container) return;

    const botRunning = !!(bot && bot.running);
    const wdRunning = !!(watchdog && watchdog.running);

    const botLabel = escapeHtml(t('bot.heading'));
    const wdLabel = escapeHtml(t('bot.watchdog'));
    const runningText = escapeHtml(t('dashboard.botRunning'));
    const stoppedText = escapeHtml(t('dashboard.botStopped'));

    container.innerHTML = `
        <div class="bot-status-card">
            <div class="bot-status-dot ${botRunning ? 'running' : 'stopped'}"></div>
            <div class="bot-status-info">
                <span class="bot-status-label">${botLabel}</span>
                <span class="bot-status-value ${botRunning ? 'running' : 'stopped'}">${botRunning ? runningText : stoppedText}</span>
            </div>
        </div>
        <div class="bot-status-card">
            <div class="bot-status-dot ${wdRunning ? 'running' : 'stopped'}"></div>
            <div class="bot-status-info">
                <span class="bot-status-label">${wdLabel}</span>
                <span class="bot-status-value ${wdRunning ? 'running' : 'stopped'}">${wdRunning ? runningText : stoppedText}</span>
            </div>
        </div>
    `;
}

async function loadBotStatusDashboard() {
    const container = document.getElementById('bot-status-cards');
    if (!container) return;

    try {
        const [botRes, wdRes] = await Promise.all([
            secureFetch('/api/bot/status'),
            secureFetch('/api/bot/watchdog/status')
        ]);

        const bot = botRes.ok ? await botRes.json().catch(() => null) : null;
        const watchdog = wdRes.ok ? await wdRes.json().catch(() => null) : null;

        if (bot) {
            renderBotStatusCards(bot, watchdog);
        } else {
            container.innerHTML = `<span class="plugin-row-description">${escapeHtml(t('dashboard.botStatusLoadError'))}</span>`;
        }
    } catch (_) {
        container.innerHTML = `<span class="plugin-row-description">${escapeHtml(t('dashboard.botStatusLoadError'))}</span>`;
    }
}

// Beispiel: Statistiken laden und Karten befüllen
function loadStats() {
    secureFetch('/api/stats')
        .then(res => res.json())
        .then(data => {
            const admins = data.admins_online ?? data.admins ?? 0;
            const members = data.members_online ?? data.members ?? 0;
            document.getElementById('card-admins').innerText = admins;
            document.getElementById('card-members').innerText = members;
        });
}

// Plugin-Status für das Dashboard laden
function loadPlugins() {
    secureFetch('/api/plugins')
        .then(res => res.json())
        .then(plugins => {
            cachedPlugins = plugins;
            renderPluginStatusList(cachedPlugins);
        });
}

function loadSettings() {
    secureFetch('/api/settings/ts3')
        .then(res => {
            if (!res.ok) {
                throw new Error('settings failed');
            }
            return res.json();
        })
        .then(data => {
            cachedTS3Settings = data;
            renderSettingsSummary(cachedTS3Settings);
        })
        .catch(() => {
            const container = document.getElementById('settings-content');
            if (container) {
                container.innerHTML = `<span class="plugin-row-description">${escapeHtml(t('dashboard.settingsLoadError'))}</span>`;
            }
        });
}

document.addEventListener('DOMContentLoaded', () => {
    loadBotStatusDashboard();
    loadStats();
    loadPlugins();
    loadSettings();
    botStatusInterval = setInterval(loadBotStatusDashboard, 10000);
});

window.addEventListener('uc-language-changed', () => {
    loadBotStatusDashboard();
    if (cachedPlugins.length) {
        renderPluginStatusList(cachedPlugins);
    }
    if (cachedTS3Settings) {
        renderSettingsSummary(cachedTS3Settings);
    }
});
