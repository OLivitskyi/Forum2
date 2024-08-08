import { loadAndRenderSinglePost, loadAndRenderComments } from "../eventHandlers.js";
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
      <form class="form" id="create-comment-form">
        <input type="hidden" id="post-id" name="post_id" value="${this.params.id}">
        <div class="form__group">
          <label for="content">Comment</label>
          <textarea id="content" name="content" required></textarea>
        </div>
        <div id="comments-container">
        <!-- Comments will be dynamically loaded here -->
        </div>
        <button type="submit">Add Comment</button>
      </form>
    `;
    return getLayoutHtml(content);
    }
    async postRender() {
      await loadAndRenderSinglePost(this.params.id);
      await loadAndRenderComments(this.params.id);
    }
}