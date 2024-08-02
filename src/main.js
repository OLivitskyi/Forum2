// Takes in loginForm or createAccountForm, type (success or error) and message defined in html 
function setFormMessage(formElement, type, message) {
    const messageElement = formElement.querySelector(".form__message");
    messageElement.textContent = message;
    messageElement.classList.remove("form__message--success", "form__message--error");
    messageElement.classList.add(`form__message--${type}`); // adds type of message provided by parameter
}
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
document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.querySelector("#login");
    const createAccountForm = document.querySelector("#createAccount");
    document.querySelector("#linkCreateAccount").addEventListener("click", e => {
        e.preventDefault(); // will not redirect to href
        loginForm.classList.add("form--hidden");
        createAccountForm.classList.remove("form--hidden");
    });
    document.querySelector("#linkLogin").addEventListener("click", e => {
        e.preventDefault();
        loginForm.classList.remove("form--hidden");
        createAccountForm.classList.add("form--hidden");
    });
    loginForm.addEventListener("submit", e => {
        e.preventDefault();
        setFormMessage(loginForm, "error", "Invalid username/password combination");
    });
    // // On every input field, checks if user takes focus off username field, and produces error message if < 10 characters in length
    // document.querySelectorAll(".form__input").forEach(inputElement => {
    //     inputElement.addEventListener("blur", e => {
    //         if (e.target.id === "signupUsername" && e.target.value.length > 0 && e.target.value.length < 10) {
    //             setInputError(inputElement, "Username must be at least 10 characters in length");
    //         }
    //     });
    //     inputElement.addEventListener("input", e => {
    //         clearInputError(inputElement); // clear errors set against input field when user clicks
    //     });
    // });
});