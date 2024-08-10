import { loadAndRenderSinglePost } from "../handlers/postHandlers.js";
import { loadAndRenderComments, handleCreateCommentFormSubmit } from "../handlers/commentHandlers.js";
import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Post");
    }

    async getHtml() {
        const content = `
            <div id="single-post-container">
                <!-- Post content will be dynamically loaded here -->
            </div>
            <div id="comments-container">
                <!-- Comments will be dynamically loaded here -->
            </div>
            <form class="form" id="create-comment-form">
                <input type="hidden" id="post-id" name="post_id" value="${this.params.id.trim()}">
                <div class="form__group">
                    <label for="content">Comment</label>
                    <textarea id="content" name="content" required></textarea>
                </div>
                <div id="comment-message" class="form__message"></div>
                <button type="submit">Add Comment</button>
            </form>
        `;
        return getLayoutHtml(content);
    }

    async postRender() {
        try {
            await loadAndRenderSinglePost(this.params.id);
            await loadAndRenderComments(this.params.id);
        } catch (error) {
            console.error("Failed to load post or comments:", error);
        }

        handleCreateCommentFormSubmit(this.params.id, () => {
            document.getElementById("content").value; 
        }, (error) => {
            console.error("Failed to submit comment:", error);
        });
    }
}