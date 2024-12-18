import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Messages");
    }

    async getHtml() {
        const content = `
            <h1>Welcome to your Messages, User!</h1>
            <div class="date"></div>
            <div class="insights"></div>
        `;
        return getLayoutHtml(content);
    }
}
