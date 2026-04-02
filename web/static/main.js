if (typeof requireAuth === 'function') {
    requireAuth();
}

// Lädt die Navigation aus nav.html und fügt sie in #nav ein
fetch('nav.html')
    .then(response => response.text())
    .then(html => {
        const nav = document.getElementById('nav');
        if (!nav) return;
        nav.innerHTML = html;

        if (typeof applyTranslations === 'function') {
            applyTranslations(nav);
        }
        if (typeof initLanguageSwitchers === 'function') {
            initLanguageSwitchers(nav);
        }

        const current = window.location.pathname || '/';
        nav.querySelectorAll('a[href]').forEach((link) => {
            const href = link.getAttribute('href');
            if (!href || href === '#' || href.startsWith('http')) return;
            if (href === '/' && current === '/') {
                link.classList.add('active-link');
                return;
            }
            if (href !== '/' && current.startsWith(href)) {
                link.classList.add('active-link');
            }
        });

        const logoutBtn = document.getElementById('nav-logout');
        if (logoutBtn && typeof logout === 'function') {
            logoutBtn.addEventListener('click', (e) => {
                e.preventDefault();
                logout();
            });
        }
    });

// Beispiel: Statistiken von der API laden
const fetchFn = typeof authFetch === 'function' ? authFetch : fetch;
fetchFn('/api/stats')
    .then(res => {
        if (!res.ok) {
            throw new Error('stats failed');
        }
        return res.json();
    })
    .then(data => {
        const admins = data.admins_online ?? data.admins ?? 0;
        const members = data.members_online ?? data.members ?? 0;
        const el = document.getElementById('stats-content');
        if (el) el.innerHTML =
            `<b>${typeof t === 'function' ? t('dashboard.supportersOnline') : 'Supporters online'}:</b> ${admins}<br><b>${typeof t === 'function' ? t('dashboard.membersOnline') : 'Members online'}:</b> ${members}`;
    })
    .catch(() => {
        const el = document.getElementById('stats-content');
        if (el) el.innerText = typeof t === 'function' ? t('dashboard.statsLoadError') : 'Failed to load statistics.';
    });
