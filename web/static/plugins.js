const secureFetch = typeof authFetch === 'function' ? authFetch : fetch;
let cachedPlugins = [];

function renderPlugins(plugins) {
    const container = document.getElementById('plugins-page-list');
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
            <label class="switch" aria-label="${pluginLabel}">
                <input type="checkbox" ${plugin.active ? 'checked' : ''}>
                <span class="slider"></span>
            </label>
        `;

        const checkbox = row.querySelector('input');
        checkbox.addEventListener('change', async (event) => {
            const previous = !event.target.checked;
            event.target.disabled = true;

            try {
                const res = await secureFetch('/api/plugins/toggle', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: plugin.name, active: event.target.checked })
                });

                if (!res.ok) {
                    throw new Error(await res.text());
                }
            } catch (error) {
                event.target.checked = previous;
                window.alert(`${t('common.error')}: ${error.message}`);
            } finally {
                event.target.disabled = false;
            }
        });

        container.appendChild(row);
    });
}

async function loadPluginsPage() {
    const container = document.getElementById('plugins-page-list');

    try {
        const res = await secureFetch('/api/plugins');
        if (!res.ok) {
            throw new Error(await res.text());
        }

        const plugins = await res.json();
        cachedPlugins = plugins;
        renderPlugins(cachedPlugins);
    } catch (error) {
        if (container) {
            container.textContent = `${t('common.error')}: ${error.message}`;
        }
    }
}

document.addEventListener('DOMContentLoaded', loadPluginsPage);
window.addEventListener('uc-language-changed', () => {
    if (cachedPlugins.length) {
        renderPlugins(cachedPlugins);
    }
});