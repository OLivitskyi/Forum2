import { sendMessage } from "../websocket.js";
import { renderUserList } from "../components/userList.js";


export const updateAndSortUserList = (users, currentUserId) => {
    renderUserList("user-status-list", users, currentUserId, (userId) => {
        console.log(`User ${userId} clicked in global user list.`);
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
