import { sendPrivateMessage } from "../websocket.js";
import { debounce } from "../routeUtils.js";

let currentReceiverID = null;
let offset = 0;
let loading = false;
let isAutoScrollEnabled = true;

export const setupMessageForm = () => {
  const sendButton = document.getElementById("send-button");
  const messageInput = document.getElementById("message-input");

  sendButton.addEventListener("click", () => {
    if (currentReceiverID && messageInput.value) {
      sendPrivateMessage(currentReceiverID, messageInput.value);
      messageInput.value = "";
    }
  });
};

export const loadMessages = async (receiverID, loadMore = false) => {
  try {
    const limit = 10;
    offset = loadMore ? offset + limit : 0;

    const response = await fetch(
      `/api/get-messages?user_id=${receiverID}&limit=${limit}&offset=${offset}`
    );
    const messages = await response.json();

    const messageList = document.querySelector(".message-list");

    if (!loadMore) {
      messageList.innerHTML = "";
    }

    if (Array.isArray(messages) && messages.length > 0) {
      messages.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));

      messages.forEach((msg) => {
        const messageElement = document.createElement("div");
        messageElement.classList.add(
          msg.sender_id === receiverID ? "other-user-message" : "user-message"
        );
        messageElement.innerHTML = `
          <div class="message-content">
            <strong>${msg.sender_name || "Unknown User"}:</strong> ${msg.content
          }
            <div class="message-time">${new Date(
            msg.created_at
          ).toLocaleString()}</div>
          </div>
        `;
        messageList.prepend(messageElement);
      });

      if (!loadMore && isAutoScrollEnabled) {
        messageList.scrollTop = messageList.scrollHeight;
      }
    } else if (!loadMore) {
      console.log("No messages found for this user.");
      messageList.innerHTML = "<p>No previous messages.</p>";
    }
  } catch (error) {
    console.error("Failed to load messages:", error);
  }
};

export const setupMessageListScroll = () => {
  const messageList = document.querySelector(".message-list");

  messageList.addEventListener(
    "scroll",
    debounce(() => {
      if (messageList.scrollTop === 0 && !loading) {
        const oldScrollHeight = messageList.scrollHeight;

        loading = true;
        loadMessages(currentReceiverID, true).finally(() => {
          loading = false;

          const newScrollHeight = messageList.scrollHeight;
          messageList.scrollTop = newScrollHeight - oldScrollHeight;
        });
      }

      isAutoScrollEnabled =
        messageList.scrollTop + messageList.clientHeight >=
        messageList.scrollHeight;
    }, 200)
  );
};

export const setCurrentReceiver = (receiverID) => {
  currentReceiverID = receiverID;
  offset = 0;
  loadMessages(receiverID);
};

export const handlePrivateMessage = (message) => {
  const messageList = document.querySelector(".message-list");
  const messagesLink = document.getElementById("messages");

  if (messageList) {
    const isAtBottom =
      messageList.scrollTop + messageList.clientHeight >=
      messageList.scrollHeight;

    const messageElement = document.createElement("div");
    messageElement.classList.add(
      message.sender_id === currentReceiverID
        ? "other-user-message"
        : "user-message",
      "new"
    );
    messageElement.innerHTML = `
      <div class="message-content">
        <strong>${message.sender_name}:</strong> ${message.content}
        <div class="message-time">${new Date(
      message.timestamp
    ).toLocaleString()}</div>
      </div>
    `;

    messageList.prepend(messageElement);

    if (isAtBottom) {
      messageList.scrollTop = messageList.scrollHeight;
    }

    setTimeout(() => {
      messageElement.classList.remove("new");
    }, 3000);
  } else {
    showPopupNotification("You have a new message!");

    const messageCountElement = messagesLink.querySelector(".message-count");
    const currentCount = parseInt(messageCountElement.textContent, 10) || 0;
    messageCountElement.textContent = currentCount + 1;
  }
};

const showPopupNotification = (message) => {
  const notification = document.getElementById("popup-notification");
  if (notification) {
    notification.textContent = message;
    notification.style.display = "block";

    setTimeout(() => {
      notification.style.display = "none";
    }, 3000);
  }
};

export const markMessagesAsRead = (receiverID) => {
  const messagesLink = document.getElementById("messages");
  const messageCountElement = messagesLink.querySelector(".message-count");
  const currentCount = parseInt(messageCountElement.textContent, 10) || 0;

  if (receiverID === currentReceiverID) {
    const unreadMessages = document.querySelectorAll(
      ".other-user-message.unread"
    );
    const unreadCount = unreadMessages.length;
    unreadMessages.forEach((msg) => msg.classList.remove("unread"));

    messageCountElement.textContent = Math.max(currentCount - unreadCount, 0);
  }
};