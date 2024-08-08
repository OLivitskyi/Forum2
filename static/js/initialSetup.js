import { navigateTo } from './routeUtils.js';
import { router } from './router.js';

document.addEventListener("DOMContentLoaded", () => {
    const handleLinkClick = (e) => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    };

    document.body.addEventListener("click", handleLinkClick);
    router();
});
