import { createCategory, getCategories } from '../api.js';
import { setFormMessage } from '../formHandler.js';

export const handleCreateCategoryFormSubmit = () => {
    const form = document.getElementById("create-category-form");
    if (form) {
        form.removeEventListener("submit", handleCreateCategorySubmit);
        form.addEventListener("submit", handleCreateCategorySubmit);
    }
};

async function handleCreateCategorySubmit(e) {
    e.preventDefault();
    const categoryName = document.getElementById("category-name").value;
    const messageElement = document.getElementById("category-message");

    if (!categoryName) {
        setFormMessage(messageElement, "error", "Category name is required");
        return;
    }

    const body = { name: categoryName };

    try {
        const response = await createCategory(body);
        if (response) {
            setFormMessage(messageElement, "success", "Category created successfully");
            document.getElementById("category-name").value = "";
            await loadCategories();
        } else {
            setFormMessage(messageElement, "error", "Failed to create category.");
        }
    } catch (error) {
        setFormMessage(messageElement, "error", "An error occurred. Please try again.");
    }
}

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
        const messageElement = document.getElementById("category-message");
        setFormMessage(messageElement, "error", error.message || "Failed to load categories");
    }
};
