import { navigateTo } from './router.js';
import { sendRequest, logout } from './auth.js';
import { showError, clearError } from './errorHandler.js';

export const handleLoginFormSubmit = (clearError, showError) => {
    const loginForm = document.getElementById("login");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            clearError();

            const formData = new FormData(e.target);
            const response = await fetch("/", {
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
};

export const handleLogout = () => {
    const logoutButton = document.getElementById("logout");
    if (logoutButton) {
        logoutButton.addEventListener("click", async (e) => {
            e.preventDefault();
            await logout();
        });
    }
};

export const handleCreatePostFormSubmit = (clearError, showError) => {
    const form = document.getElementById("create-post-form"); // Ensure this ID matches the form in CreatePost.js
    if (form) {
        form.addEventListener("submit", async (e) => {
            e.preventDefault();
            clearError();

            const title = document.getElementById("title").value;
            const content = document.getElementById("content").value;
            const categories = Array.from(document.querySelectorAll(".pill--selected")).map(pill => pill.innerText);

            if (!title || !content || categories.length === 0) {
                showError("Title, content and at least one category are required");
                return;
            }

            const body = {
                title,
                content,
                category_ids: categories
            };

            try {
                const response = await sendRequest("/api/create-post", "POST", body);

                if (response.ok) {
                    clearError();
                    navigateTo("/homepage");
                } else {
                    const errorText = await response.text();
                    showError(errorText || "Failed to create post");
                }
            } catch (error) {
                showError("An error occurred. Please try again.");
            }
        });
    }
};
