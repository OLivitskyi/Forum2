import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { handleCreateCategoryFormSubmit, loadCategories } from "../eventHandlers.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Create Category");
    }

    async getHtml() {
        const content = `
            <form class="form" id="create-category-form">
                <h1>Create a Category</h1>
                <div class="create-category-container">
                    <input type="text" id="category-name" name="category-name" class="form__input" placeholder="Category name">
                    <button class="pill pill-submit" id="create-category-button" type="submit">CREATE</button>
                    <div id="category-message" class="form__message"></div>
                </div>
            </form>
        `;
        return getLayoutHtml(content);
    }

    async postRender() {
        await loadCategories();
        const form = document.getElementById("create-category-form");
        if (form) {
            handleCreateCategoryFormSubmit(form);
        }
    }
}
