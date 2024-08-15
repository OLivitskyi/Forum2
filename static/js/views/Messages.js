import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { requestUserStatus } from "../handlers/userStatusHandlers.js";
import {
    setupMessageForm,
    loadMessages,
    setupMessageListScroll,
    setCurrentReceiver,
} from "../handlers/messageHandlers.js";
import { getUserInfo } from "../api.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Messages");
    }

    async getHtml() {
        const content = `
            <h1>Welcome to your Messages, User!</h1>
            <div class="message-container">
                <div class="box" id="box1">
                    <!-- User list will be dynamically populated here -->
                </div>
                <div class="box" id="box2">
                    <div class="message-list">
                        <!-- Messages will be dynamically populated here -->
                    </div>
                    <div class="chatbox">
                        <input type="text" id="message-input" placeholder="Write your message here">
                        <button id="send-button">SEND</button>
                    </div>
                </div>
            </div>
        `;

        setTimeout(async () => {
            this.currentUserInfo = await getUserInfo();
            this.initializeEvents();
            requestUserStatus();
        }, 0);

        return getLayoutHtml(content);
    }

    initializeEvents() {
        const userList = document.getElementById("box1");
        const currentUserId = this.currentUserInfo.user_id;

        userList.addEventListener("click", (event) => {
            const target = event.target.closest(".user-box");
            if (target) {
                const receiverID = target.dataset.userId;

                if (receiverID === currentUserId) {
                    alert("You cannot send messages to yourself.");
                    return;
                }

                setCurrentReceiver(receiverID);
                loadMessages(receiverID);
            }
        });

        setupMessageForm();
        setupMessageListScroll();
    }
}