import { renderUserList } from "../components/userList.js";
import { requestUserStatus } from "../handlers/userStatusHandlers.js";
import { getUserInfo } from "../api.js";

export const getLayoutHtml = (content) => {
  const layout = `
        <div class="container">
            <aside>
                <div class="top">
                    <div class="logo">
                        <span class="material-icons-sharp logo-icon">app_registration</span>
                        <h2>FOR<span class="danger">UM</span></h2>
                    </div>
                    <div class="close" id="close-btn">
                        <span class="material-icons-sharp">close</span>
                    </div>
                </div>
                <div class="sidebar">
                    <a href="homepage" data-link id="homepage">
                        <span class="material-icons-sharp">grid_view</span>
                        <h3>Dashboard</h3>
                    </a>
                    <a href="messages" data-link id="messages">
                        <span class="material-icons-sharp">mail_outline</span>
                        <h3>Messages</h3>
                        <span class="message-count">0</span>
                    </a>
                    <a href="create-post" data-link id="create-post">
                        <span class="material-icons-sharp">add</span>
                        <h3>Create Post</h3>
                    </a>
                    <a href="create-category" data-link id="create-category">
                        <span class="material-icons-sharp">category</span>
                        <h3>Create Category</h3>
                    </a>
                    <a href="#" id="logout">
                        <span class="material-icons-sharp">logout</span>
                        <h3>Logout</h3>
                    </a>
                </div>
            </aside>
            <main>
                ${content}
            </main>
            <div class="user-status-container">
                <div id="user-status-list" class="user-status-list"></div>
            </div>
        </div>
        <div id="popup-notification">You have a new message!</div>
    `;

  setTimeout(async () => {
    try {
      const currentUserInfo = await getUserInfo();
      if (!currentUserInfo || !currentUserInfo.user_id) {
        console.error("Failed to retrieve user info");
        return;
      }

      const currentUserId = currentUserInfo.user_id;
      requestUserStatus();

      const users = JSON.parse(localStorage.getItem("users")) || [];
      if (!users.length) {
        console.warn("No users found in localStorage");
      }

      const userStatusListElement = document.getElementById("user-status-list");
      if (!userStatusListElement) {
        console.error("User status list element not found");
        return;
      }

      renderUserList("user-status-list", users, currentUserId, (userId) => {
        console.log(`User ${userId} clicked in global user list.`);
      });
    } catch (error) {
      console.error("An error occurred while setting up the user list:", error);
    }
  }, 0);

  return layout;
};
