import Homepage from "./views/Homepage.js";
import Login from "./views/Login.js";
// import Posts from "./views/Posts.js";
import Registration from "./views/Registration.js";
// import Settings from "./views/Settings.js";

const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));
};

const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

const isAuthenticated = async () => {
    const response = await fetch("/validate-session", {
        method: "GET"
    });
    return response.status === 200;
}

const showError = (message) => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = message;
    }
};

const clearError = () => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = '';
    }
};

const router = async () => {
    const routes = [
        { path: "/", view: Login }, // class reference
        // { path: "/posts", view: Posts },
        { path: "/registration", view: Registration },
        { path: "/homepage", view: Homepage },
        { path: "/logout", view: Login }
        // { path: "/settings", view: Settings }
    ];
    // Test each route for potential match
    const potentialMatches = routes.map(route => {
        return {
            
            route: route,
            result: location.pathname.match(pathToRegex(route.path))
            // isMatch: location.pathname === route.path
        };
    });

    let match = potentialMatches.find(potentialMatch => potentialMatch.result!== null);
    if (!match) {
        match = {
            route: routes[0], // set default as "/"
            result: [location.pathname]
        };
    }
    if (match.route.protected) {
        const auth = await isAuthenticated();
        if (!auth) {
            navigateTo("/");
            return
        }
    }

    const view = new match.route.view(getParams(match)); // creates new instance of view at match route

    document.querySelector("#app").innerHTML = await view.getHtml(); // select app element

    // Apples styling to input errors and adds message
function setInputError(inputElement, message) {
    inputElement.classList.add("form__input--error");
    inputElement.parentElement.querySelector(".form__input-error-message").textContent = message;
}

// Removes error if user fixes it according to conditions (undo setInputError)
function clearInputError(inputElement) {
    inputElement.classList.remove("form__input--error");
    inputElement.parentElement.querySelector(".form__input-error-message").textContent = "";
}

    // LOGOUT
    let logout = document.getElementById("logout");
    if (logout) {
        logout.addEventListener("click", async e => {
            e.preventDefault()
            navigateTo(e.target.href);
        })
    }

    document.addEventListener("DOMContentLoaded", () => {
        document.body.addEventListener("click", e => {
            if (e.target.matches("[data-link]")) {
                e.preventDefault();
                if (e.target.id == "logout") {
                    logoutUser(e);
                } else {
                    navigateTo(e.target.href);
                }
            }
        });
        router();
    });
    
    const logoutUser = async (e) => {
        e.preventDefault();
        let response = await fetch("/logout", {
            method: "GET",
        });
        if (response.ok) {
            navigateTo("/");
        }
    };

  
    // LOGIN
    let login = document.getElementById("login");
    if (login) {
        login.addEventListener("submit", async e => {
            e.preventDefault()
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
                showError(errorText || "Login failed. Please try again.")
            }
            
        })
    }


    // REGISTRATION
        document.querySelectorAll(".form__input").forEach(inputElement => {
        inputElement.addEventListener("blur", e => {
            if (e.target.id === "signupUsername" && e.target.value.length > 0 && e.target.value.length < 10) {
                setInputError(inputElement, "Username must be at least 10 characters in length");
            }
        });

        inputElement.addEventListener("input", e => {
            clearInputError(inputElement); // clear errors set against input field when user clicks
        });
    });

    let createAccount = document.getElementById("createAccount");
    if (createAccount) {
        console.log(createAccount);
        createAccount.addEventListener("submit", async e => {
            e.preventDefault();
            console.log("submit clicked");
            
            // console.log(e.target.getElementById("signupUsername")
            const formData = new FormData(e.target);
            try {
                let response = await fetch("/registration", {
                    method: "POST",
                    body: formData,
                })
                console.log(formData);
                var responseText = await response.text() 
                console.log(responseText);
                
                console.log(response);
                if (response.ok) {
                    navigateTo("/");
                } else {
                    const errorText = await response.text();
                    showError(errorText || "Registration failed. Please try again");
                }
            }
            catch (error) {
                console.log(error);
                showError("An unexpected error occurred. Please try again.")
            }
            
            
        })
    }
};

window.addEventListener("popstate", router); // upon user navigating history, run router 

// Reroutes to selected route (upon click) according to data-link attribute without reloading page 
document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            if (e.target.id == "submit")
            {
                // Submit form data to server

                console.log(e.target.form);
            }
            console.log(e.target);
            navigateTo(e.target.href);
        }
    })
    router();
});
window.onload