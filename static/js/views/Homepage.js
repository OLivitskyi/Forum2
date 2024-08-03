import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Homepage");
    }

    async getHtml() {
        const content = `
            <h1>Have something to share?</h1>
            <div class="create-post-container">
                <h2>What's on your mind?</h2>
                <div class="container-post">
                    <input type="text" name="post-title" class="post-title" placeholder="Title of your post">
                    <input type="text" name="post-content" class="post-subject" placeholder="Write a post">
                </div>
                <button class="pill" type="button">General</button>
                <button class="pill pill--selected" type="button">Travel</button>
                <button class="pill" type="button">Hobbies</button>
                <button class="pill" type="button">Gaming</button>
                <button class="pill pill-submit" type="button">POST</button>
            </div>
        `;
        return getLayoutHtml(content);
    }
}
