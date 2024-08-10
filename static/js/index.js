import { router } from './router.js';
import { navigateTo } from './routeUtils.js';

document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
        e.preventDefault(); 
        navigateTo(e.target.href);
    }
});