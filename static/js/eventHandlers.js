import { navigateTo } from './router.js';
import { sendRequest, logout } from './auth.js';
import { setFormMessage } from './formHandler.js';
import { getCategories } from './api.js';
import { showError, clearError } from './errorHandler.js';

export const handleLoginFormSubmit = (clearError, showError) => {
    const loginForm = document.getElementById("login");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            clearError();

            const formData = new FormData(e.target);
            const response = await fetch("/api/login", {
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
    const form = document.getElementById("create-post-form");
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

export const setupCreateCategoryForm = () => {
    const createCategoryButton = document.getElementById("create-category-button");
    if (createCategoryButton) {
        createCategoryButton.addEventListener("click", async () => {
            const categoryName = document.getElementById("category-name").value;
            const messageElement = document.getElementById("category-message");
            if (categoryName === "") {
                setFormMessage(messageElement, "error", "Category name is required");
                return;
            }
            const response = await fetch("/api/create-category", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                },
                body: `name=${categoryName}`
            });
            if (response.ok) {
                setFormMessage(messageElement, "success", "Category created successfully");
                document.getElementById("category-name").value = "";
            } else {
                const errorText = await response.text();
                setFormMessage(messageElement, "error", `Error: ${errorText}`);
            }
        });
    }
};

export const loadCategories = async () => {
    try {
        const categories = await getCategories();
        const categorySelect = document.getElementById("category-select");
        if (!categorySelect) return;

        categories.forEach(category => {
            const option = document.createElement("option");
            option.value = category.id;
            option.textContent = category.name;
            categorySelect.appendChild(option);
        });
    } catch (error) {
        console.error("Failed to load categories:", error);
        const messageElement = document.getElementById("category-message");
        if (messageElement) {
            setFormMessage(messageElement, "error", error.message || "Failed to load categories");
        }
    }
};
