import { createCategory, getCategories } from '../api.js';
import { showError, clearError } from '../errorHandler.js';
import { setFormMessage } from '../formHandler.js';

export const handleCreateCategoryFormSubmit = (clearError, showError) => {
    const form = document.getElementById("create-category-form");
    if (form) {
        form.removeEventListener("submit", handleCreateCategorySubmit);
        form.addEventListener("submit", handleCreateCategorySubmit, { once: true });
    }
};

async function handleCreateCategorySubmit(e) {
    e.preventDefault();
    clearError();
    const categoryName = document.getElementById("category-name").value;
    const messageElement = document.getElementById("category-message");
    if (!categoryName) {
        showError("Category name is required");
        return;
    }
    const body = {
        name: categoryName
    };
    try {
        const response = await createCategory(body);
        if (response.ok) {
            setFormMessage(messageElement, "success", "Category created successfully");
            document.getElementById("category-name").value = "";
            await loadCategories();
        } else {
            const errorText = await response.text();
            setFormMessage(messageElement, "error", `Error: ${errorText}`);
        }
    } catch (error) {
        showError("An error occurred. Please try again.");
    }
};

export const loadCategories = async () => {
    try {
        const categorySelect = document.getElementById("category-select");
        if (!categorySelect) return;
        categorySelect.innerHTML = '<option value="">Select a category</option>';
        const categories = await getCategories();
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
