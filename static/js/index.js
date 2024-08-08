import { router } from './router.js';
import { navigateTo } from './routeUtils.js'; // Оновлено шлях до routeUtils.js

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.prevent.preventDefault();
            navigateTo(e.target.href);
        }
    });
    router();
});
