import Homepage from "./views/Homepage.js";
import Login from "./views/Login.js";
import Registration from "./views/Registration.js";
import CreatePost from "./views/CreatePost.js";
import Messages from "./views/Messages.js";
import CreatePostCategory from "./views/CreatePostCategory.js";
import { isAuthenticated } from './auth.js';
import { handleLoginFormSubmit, handleLogout, handleCreatePostFormSubmit } from './eventHandlers.js';
import { showError, clearError } from './errorHandler.js';
import { setInputError, clearInputError, setupFormSwitching, setupFormValidation } from './formHandler.js';

// Function to convert path to regex for route matching
const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

// Function to get parameters from the matched route
const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);
    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));
};

// Function to navigate to a new URL and call the router
export const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

// Main router function to handle route changes
export const router = async () => {
    const routes = [
        { path: "/", view: Login },
        { path: "/registration", view: Registration },
        { path: "/homepage", view: Homepage, protected: true },
        { path: "/logout", view: Login },
        { path: "/create-post", view: CreatePost, protected: true },
        { path: "/messages", view: Messages, protected: true },
        { path: "/create-category", view: CreatePostCategory, protected: true }
    ];
    
    // Find the matching route
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

    // Check if the route is protected and if the user is authenticated
    if (match.route.protected) {
        const auth = await isAuthenticated();
        if (!auth) {
            navigateTo("/");
            return;
        }
    }

    // Load the view for the matched route
    const view = new match.route.view(getParams(match));
    document.querySelector("#app").innerHTML = await view.getHtml();
    if (view.postRender) {
        view.postRender();
    }

    // Call setup functions after view is loaded
    handleLoginFormSubmit(clearError, showError);
    handleLogout();
    setupFormSwitching();
    setupFormValidation();
    handleCreatePostFormSubmit(clearError, showError);

    // Toggle category selection for post creation
    document.querySelectorAll(".pill").forEach(pill => {
        pill.addEventListener("click", () => pill.classList.toggle("pill--selected"));
    });

    // Messages event handler
    let messages = document.getElementById("messages");
    if (messages) {
        messages.addEventListener("click", async e => {
            console.log("messages clicked");
            e.preventDefault();
            navigateTo("/messages");
        });
    }

    // Homepage event handler
    let homepage = document.getElementById("homepage");
    if (homepage) {
        homepage.addEventListener("click", async e => {
            console.log("homepage clicked");
            e.preventDefault();
            navigateTo("/homepage");
        });
    }

    // Input validation for registration form
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

    // Registration event handler
    let createAccount = document.getElementById("createAccount");
    if (createAccount) {
        createAccount.addEventListener("submit", async e => {
            e.preventDefault();
            const formData = new FormData(e.target);
            let response = await fetch("/registration", {
                method: "POST",
                body: formData,
            });
            if (response.ok) {
                navigateTo("/");
            } else {
                const errorText = await response.text();
                showError(errorText || "Registration failed. Please try again.");
            }
        });
    }
};

// Handle back/forward button navigation
window.addEventListener("popstate", router);

// Initial setup and event listeners
document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });
    router();
});
