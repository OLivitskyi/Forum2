import AbstractView from "./AbstractView.js";
import { getLayoutHtml } from "./layout.js";
import { requestUserStatus } from "../websocket.js"; // Імпортуємо нову функцію

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
            <div class="message user-message">
              <div class="message-content">This is a message from the user</div>
            </div>
            <div class="message other-user-message">
              <div class="message-content">This is a message from another user</div>
            </div>
          </div>
          <div class="chatbox">
            <input type="text" id="message-input" placeholder="Write your message here">
            <button id="send-button">SEND</button>
          </div>
        </div>
      </div>
    `;

    // Викликаємо функцію для запиту статусів користувачів
    requestUserStatus();

    return getLayoutHtml(content);
  }
}