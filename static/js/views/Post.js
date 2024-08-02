import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Post");
  }

  async getHtml() {
    return `
      <div id="post-container">
        <!-- Post content will be dynamically loaded here -->
      </div>
      <form class="form" id="create-comment-form">
        <input type="hidden" id="post-id" name="post_id" value="${this.params.id}">
        <div class="form__group">
          <label for="content">Comment</label>
          <textarea id="content" name="content" required></textarea>
        </div>
        <button type="submit">Add Comment</button>
      </form>
    `;
  }
}
