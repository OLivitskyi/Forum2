import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { loadAndRenderPosts } from "../handlers/postHandlers.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Homepage");
    }

    async getHtml() {
        const content = `
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
