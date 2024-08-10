import { routes } from './routes.js';
import { pathToRegex, getParams, navigateTo } from './routeUtils.js';
import { isAuthenticated } from './auth.js';
import { handleLoginFormSubmit, handleLogout, handleCreatePostFormSubmit, handleCreateCategoryFormSubmit } from './eventHandlers.js';
import { showError, clearError } from './errorHandler.js';
import { setInputError, clearInputError, setupFormSwitching, setupFormValidation } from './formHandler.js';

export const router = async () => {
    console.log("Routing started");

    const potentialMatches = routes.map(route => {
        return {
            route: route,
            result: location.pathname.match(pathToRegex(route.path))
        };
    });

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);
    if (!match) {
        match = {
            route: routes[0],
            result: [location.pathname]
        };
    }

    if (match.route.protected) {
        const auth = await isAuthenticated();
        if (!auth) {
            console.log("User not authenticated, redirecting to login");
            navigateTo("/");
            return;
        }
    }

    const view = new match.route.view(getParams(match));
    document.querySelector("#app").innerHTML = await view.getHtml();
    if (view.postRender) {
        view.postRender();
    }

    console.log(`View for ${location.pathname} loaded`);

    handleLoginFormSubmit();
    handleLogout();
    setupFormSwitching();
    setupFormValidation();
    handleCreatePostFormSubmit(clearError, showError);
    handleCreateCategoryFormSubmit(clearError, showError);
    setupElementHandlers();
};

const setupElementHandlers = () => {
    document.querySelectorAll(".pill").forEach(pill => {
        pill.addEventListener("click", () => pill.classList.toggle("pill--selected"));
    });

    const messages = document.getElementById("messages");
    if (messages) {
        messages.addEventListener("click", async e => {
            console.log("messages clicked");
            e.preventDefault();
            navigateTo("/messages");
        });
    }

    const createpost = document.getElementById("create-post");
    if (createpost) {
        createpost.addEventListener("click", async e => {
            console.log("create post clicked");
            e.preventDefault();
            navigateTo("/create-post");
        });
    }

    const createcategory = document.getElementById("create-category");
    if (createcategory) {
        createcategory.addEventListener("click", async e => {
            console.log("create category clicked");
            e.preventDefault();
            navigateTo("/create-category");
        });
    }

    if (location.pathname === "/create-category") {
        handleCreateCategoryFormSubmit(clearError, showError);
    }

    const homepage = document.getElementById("homepage");
    if (homepage) {
        homepage.addEventListener("click", async e => {
            console.log("homepage clicked");
            e.preventDefault();
            navigateTo("/homepage");
        });
    }

    document.querySelectorAll(".form__input").forEach(inputElement => {
        inputElement.addEventListener("blur", e => {
            if (e.target.id === "signupUsername") {
                if (e.target.value.length > 0 && e.target.value.length < 1) {
                    setInputError(inputElement, "Username must be at least 1 character in length");
                }
                if (e.target.value.includes("@")) {
                    setInputError(inputElement, "Username cannot include '@'");
                }
            }
        });
        inputElement.addEventListener("input", e => {
            clearInputError(inputElement);
        });
    });
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });
    router();
});
