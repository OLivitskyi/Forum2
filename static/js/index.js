import { router } from './router.js';
import { navigateTo } from './routeUtils.js'; // Оновлено шлях до routeUtils.js

document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
        e.preventDefault();  // Виправлено
        navigateTo(e.target.href);
    }
});