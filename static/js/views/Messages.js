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
import { renderUserList } from "../components/userList.js";

export default class extends AbstractView {
    constructor(params) {
      super(params);
      this.setTitle("Messages");
    }
  
    async getHtml() {
      const content = `
        <h1 id="welcome-message">Welcome to your Messages, User!</h1>
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
  
      return getLayoutHtml(content);
    }
  
    async postRender() {
      try {
        const currentUserInfo = await getUserInfo();
        const currentUserId = currentUserInfo.user_id;
        const userName = currentUserInfo.username; 
  
        const welcomeMessage = document.getElementById('welcome-message');
        if (welcomeMessage) {
          welcomeMessage.textContent = `Welcome to your Messages, ${userName}!`;
        }
  
        requestUserStatus();
        renderUserList(
          "box1",
          JSON.parse(localStorage.getItem("users")) || [],
          currentUserId,
          (receiverID) => {
            if (receiverID !== currentUserId) {
              if (document.querySelector(".message-list")) {
                setCurrentReceiver(receiverID);
              } else {
                console.error("Message list not found on the page.");
              }
            }
          }
        );
  
        setupMessageForm();
        setupMessageListScroll();
      } catch (error) {
        console.error("An error occurred during postRender:", error);
      }
    }
  }
  
