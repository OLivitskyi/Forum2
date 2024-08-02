import Homepage from "./views/Homepage.js";
import Login from "./views/Login.js";
import Registration from "./views/Registration.js";
import CreatePost from "./views/CreatePost.js";
import Messages from "./views/Messages.js";

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
const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

// Function to send an HTTP request with session token
const sendRequest = async (url, method, body = null) => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    
    // Get session token from cookies
    const sessionToken = document.cookie.split('; ').find(row => row.startsWith('session_token='));
    if (sessionToken) {
        headers.append('Authorization', `Bearer ${sessionToken.split('=')[1]}`);
    }
    
    const response = await fetch(url, {
        method: method,
        headers: headers,
        body: body ? JSON.stringify(body) : null
    });
    
    // If unauthorized, redirect to login page
    if (response.status === 401) {
        navigateTo("/");
    }
    
    return response;
};

// Function to check if the user is authenticated
const isAuthenticated = async () => {
    const response = await sendRequest("/validate-session", "GET");
    return response.status === 200;
};

// Function to display an error message
const showError = (message) => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = message;
    }
};

// Function to clear the error message
const clearError = () => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = '';
    }
};

// Main router function to handle route changes
const router = async () => {
    const routes = [
        { path: "/", view: Login },
        { path: "/registration", view: Registration },
        { path: "/homepage", view: Homepage, protected: true },
        { path: "/logout", view: Login },
        { path: "/create-post", view: CreatePost, protected: true },
        { path: "/messages", view: Messages, protected: true }
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

    // Function to set input error message
    function setInputError(inputElement, message) {
        inputElement.classList.add("form__input--error");
        inputElement.parentElement.querySelector(".form__input-error-message").textContent = message;
    }

    // Function to clear input error message
    function clearInputError(inputElement) {
        inputElement.classList.remove("form__input--error");
        inputElement.parentElement.querySelector(".form__input-error-message").textContent = "";
    }

    // Logout event handler
    let logout = document.getElementById("logout");
    if (logout) {
        logout.addEventListener("click", async e => {
            e.preventDefault();
            let response = await sendRequest("/logout", "GET");
            if (response.ok) {
                navigateTo("/");
            }
        });
    }

    // Login event handler
    let login = document.getElementById("login");
    if (login) {
        login.addEventListener("submit", async e => {
            e.preventDefault();
            const formData = new FormData(e.target);
            let response = await fetch("/", {
                method: "POST",
                body: formData,
            });
            if (response.ok) {
                clearError();
                navigateTo("/homepage");
            } else {
                const errorText = await response.text();
                showError(errorText || "Login failed. Please try again.");
            }
        });
    }
    
    // Create post event handler
    let createPost = document.getElementById("create-post");
    if (createPost) {
        createPost.addEventListener("click", async e => {
            console.log("create post clicked");
            e.preventDefault();
            navigateTo("/create-post");
        });
    }

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
