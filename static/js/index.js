import { router } from './router.js';
import { navigateTo } from './routeUtils.js';
import { initializeWebSocket } from './websocket.js';


document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
        e.preventDefault();
        navigateTo(e.target.href);
    }
});

document.addEventListener("DOMContentLoaded", () => {
    const sessionToken = localStorage.getItem('session_token');
    if (sessionToken) {
        initializeWebSocket(sessionToken);
    }
    router();
});