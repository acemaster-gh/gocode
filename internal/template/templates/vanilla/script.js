// PROJECT_NAME — gocode v2.0.0
// Vanilla JS entry point

console.log('%c PROJECT_NAME ', 
  'background:#6366f1;color:#fff;font-weight:700;padding:2px 6px;border-radius:4px');

// ── DOM ready ─────────────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
  initApp();
});

function initApp() {
  bindEvents();
}

function bindEvents() {
  const btn = document.getElementById('js-btn');
  if (btn) {
    btn.addEventListener('click', () => {
      console.log('PROJECT_NAME ready ⚡');
    });
  }
}
