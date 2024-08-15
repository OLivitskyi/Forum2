import { sendMessage } from "../websocket.js";

export const updateAndSortUserList = (users, currentUserId) => {
    const userContainer = document.getElementById("box1");
    if (!userContainer || !users.length) return;

    users.sort((a, b) => {
        const aLastMessageTime = new Date(a.last_message_time || 0);
        const bLastMessageTime = new Date(b.last_message_time || 0);

        if (aLastMessageTime > bLastMessageTime) return -1;
        if (aLastMessageTime < bLastMessageTime) return 1;

        return a.username.localeCompare(b.username);
    });

    userContainer.innerHTML = "";

    users
        .filter(user => user.user_id !== currentUserId)
        .forEach(user => {
            const userElement = document.createElement("div");
            userElement.classList.add("user-box");
            userElement.dataset.userId = user.user_id;

            const statusClass = user.is_online ? "logged-in" : "logged-out";
            userElement.innerHTML = `<span class="${statusClass}">●</span>${user.username}`;
            userContainer.appendChild(userElement);
        });
};

export const handleUserStatus = (users) => {
    console.log("Handling user status update", users);

    const currentUserId = localStorage.getItem("user_id");
    const storedUsers = JSON.parse(localStorage.getItem("users")) || [];

    users.forEach(serverUser => {
        const storedUser = storedUsers.find(user => user.user_id === serverUser.user_id);
        if (storedUser) {
            storedUser.is_online = serverUser.is_online;
        } else {
            storedUsers.push({
                ...serverUser,
                last_message_time: null
            });
        }
    });

    localStorage.setItem("users", JSON.stringify(storedUsers));

    updateAndSortUserList(storedUsers, currentUserId);
};

export const requestUserStatus = () => {
    sendMessage({ type: "request_user_status" });
};
