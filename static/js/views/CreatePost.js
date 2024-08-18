import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { handleCreatePostFormSubmit } from "../handlers/postHandlers.js";
import { loadCategories } from "../handlers/categoryHandlers.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Create Post");
  }

  async getHtml() {
    const content = `
      <form class="form" id="create-post-form">
        <h1>Create a post</h1>
        <div class="date"></div>
        <div class="insights"></div>
        <div class="create-post-container">
          <input type="text" id="title" name="post-title" class="post-title" placeholder="Title of your post">
          <textarea id="content" name="post-content" class="post-subject" placeholder="Write a post"></textarea>
          <div class="categories">
            <select id="category-select" class="form__input">
              <option value="">Select a category</option>
            </select>
            <div class="error" id="category-error" style="color: red; display: none;">Category is required.</div>
          </div>
          <button class="pill pill-submit" type="submit">POST</button>
        </div>
      </form>
    `;
    return getLayoutHtml(content);
  }

  async postRender() {
    console.log("Running postRender in CreatePost.js");
    await loadCategories();

    handleCreatePostFormSubmit();
  }
}
