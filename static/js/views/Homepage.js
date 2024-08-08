import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { loadAndRenderPosts } from "../handlers/postHandlers.js";
import { connectWebSocket, setupWebSocketHandlers } from "../websocket.js";

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
        setupWebSocketHandlers();

        console.log("Cookies:", document.cookie);
        const sessionToken = document.cookie.split('; ').find(row => row.startsWith('session_token='));
        if (sessionToken) {
            const tokenValue = sessionToken.split('=')[1];
            console.log("Session Token:", tokenValue);
            connectWebSocket(tokenValue);
        } else {
            console.error("Session token not found");
        }
    }
}
