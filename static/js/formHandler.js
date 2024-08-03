// Function to set a form message (success or error)
export const setFormMessage = (messageElement, type, message) => {
    if (messageElement) {
        messageElement.textContent = message;
        messageElement.classList.remove("form__message--success", "form__message--error");
        messageElement.classList.add(`form__message--${type}`);
    }
};

// Function to set input error message
export const setInputError = (inputElement, message) => {
    inputElement.classList.add("form__input--error");
    const errorMessageElement = inputElement.parentElement.querySelector(".form__input-error-message");
    if (errorMessageElement) {
        errorMessageElement.textContent = message;
    }
};

// Function to clear input error message
export const clearInputError = (inputElement) => {
    inputElement.classList.remove("form__input--error");
    const errorMessageElement = inputElement.parentElement.querySelector(".form__input-error-message");
    if (errorMessageElement) {
        errorMessageElement.textContent = "";
    }
};

// Function to handle switching between login and create account forms
export const setupFormSwitching = () => {
    const loginForm = document.querySelector("#login");
    const createAccountForm = document.querySelector("#createAccount");

    if (loginForm && createAccountForm) {
        document.querySelector("#linkCreateAccount").addEventListener("click", e => {
            e.preventDefault();
            loginForm.classList.add("form--hidden");
            createAccountForm.classList.remove("form--hidden");
        });

        document.querySelector("#linkLogin").addEventListener("click", e => {
            e.preventDefault();
            loginForm.classList.remove("form--hidden");
            createAccountForm.classList.add("form--hidden");
        });
    }
};

// Function to setup form validation
export const setupFormValidation = () => {
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

    // Adding event listener for category-select
    const categorySelect = document.getElementById("category-select");
    if (categorySelect) {
        categorySelect.addEventListener("change", e => {
            clearInputError(categorySelect);
        });
    }
};
