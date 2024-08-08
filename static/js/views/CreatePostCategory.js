import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { handleCreateCategoryFormSubmit, loadCategories } from "../handlers/categoryHandlers.js";
import { showError, clearError } from "../errorHandler.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Create Category");
    }

    async getHtml() {
        const content = `
            <form class="form" id="create-category-form">
                <h1>Create a Category</h1>
                <div class= "create-post-container">
                    <div class="create-category-container">
                        <input type="text" id="category-name" name="category-name" class="form__input_category" placeholder="Category name">
                        <button class="pill pill-submit" id="create-category-button" type="submit">CREATE</button>
                        <div id="category-message" class="form__message"></div>
                    </div>
                </div>
            </form>
        `;
        return getLayoutHtml(content);
    }

    async postRender() {
        await loadCategories();
        const form = document.getElementById("create-category-form");
        if (form) {
            form.removeEventListener("submit", this.handleCategoryFormSubmit);
            form.addEventListener("submit", this.handleCategoryFormSubmit);
        }
    }

    handleCategoryFormSubmit = (e) => {
        e.preventDefault();
        clearError();
        handleCreateCategoryFormSubmit(clearError, showError);
    }
}
