// =============================================================================
// script.js — PROJECT_NAME
// =============================================================================

// ── DOM Objects ───────────────────────────────────────────────────────────────
const DOM = {
    app:         document.getElementById('app'),
    themeToggle: document.getElementById('themeToggle'),
};

// ── State ─────────────────────────────────────────────────────────────────────
const state = {
    theme: localStorage.getItem('theme') || 'light',
};

// ── Functions ─────────────────────────────────────────────────────────────────
function setTheme(theme) {
    state.theme = theme;
    document.documentElement.setAttribute('data-theme', theme);
    DOM.themeToggle.querySelector('.theme-toggle__icon').textContent =
        theme === 'dark' ? '☀️' : '🌙';
    localStorage.setItem('theme', theme);
}

function toggleTheme() {
    setTheme(state.theme === 'dark' ? 'light' : 'dark');
}

function init() {
    setTheme(state.theme);
    console.log('PROJECT_NAME — ready.');
}

// ── Event Listeners ───────────────────────────────────────────────────────────
DOM.themeToggle.addEventListener('click', toggleTheme);

// ── Init ──────────────────────────────────────────────────────────────────────
init();
