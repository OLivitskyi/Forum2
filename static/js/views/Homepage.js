import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { loadAndRenderPosts } from "../eventHandlers.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Homepage");
    }

    async getHtml() {
        const content = `
            <h1>All posts in forum</h1>
            <div id="posts-container" class="posts-container">
                <!-- Posts will be inserted here -->
            </div>
        `;
        return getLayoutHtml(content);
    }

    async postRender() {
        await loadAndRenderPosts();
    }
}
