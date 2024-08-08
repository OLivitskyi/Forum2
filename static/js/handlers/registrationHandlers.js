import { showError, clearError } from "../errorHandler.js";
import { navigateTo } from "../routeUtils.js";

export const handleRegistrationFormSubmit = () => {
    const createAccountForm = document.getElementById("createAccount");
    if (createAccountForm) {
        createAccountForm.removeEventListener("submit", handleRegistrationSubmit);
        createAccountForm.addEventListener("submit", handleRegistrationSubmit);
    }
};

async function handleRegistrationSubmit(e) {
    e.preventDefault();
    clearError();
    const formData = new FormData(e.target);
    try {
        const response = await fetch("/registration", {
            method: "POST",
            body: formData,
        });
        if (response.ok) {
            navigateTo("/");
        } else {
            const errorText = await response.text();
            showError(errorText || "Registration failed. Please try again.");
        }
    } catch (error) {
        showError("An error occurred. Please try again.");
    }
}
