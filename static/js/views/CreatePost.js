import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { handleCreatePostFormSubmit } from "../eventHandlers.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("CreatePost");
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
                        <button class="pill" type="button">General</button>
                        <button class="pill pill--selected" type="button">Travel</button>
                        <button class="pill" type="button">Hobbies</button>
                        <button class="pill" type="button">Gaming</button>
                    </div>
                    <button class="pill pill-submit" type="submit">POST</button>
                </div>
            </form>
        `;
        return getLayoutHtml(content);
    }

    async postRender() {
        handleCreatePostFormSubmit();
    }
}
