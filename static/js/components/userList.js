export const renderUserList = (userContainerId, users, currentUserId, onUserClick) => {
    const userContainer = document.getElementById(userContainerId);
    if (!userContainer || !users.length) return;

    // Очищуємо контейнер перед додаванням нових користувачів
    userContainer.innerHTML = "";

    const addedUserIds = new Set();

    // Сортування і рендеринг списку користувачів
    users.sort((a, b) => {
        const aLastMessageTime = new Date(a.last_message_time || 0);
        const bLastMessageTime = new Date(b.last_message_time || 0);

        if (aLastMessageTime > bLastMessageTime) return -1;
        if (aLastMessageTime < bLastMessageTime) return 1;

        return a.username.localeCompare(b.username);
    });

    users
        .filter(user => user.user_id !== currentUserId)
        .forEach(user => {
            if (!addedUserIds.has(user.user_id)) {
                const userElement = document.createElement("div");
                userElement.classList.add("user-box");
                userElement.dataset.userId = user.user_id;

                const statusClass = user.is_online ? "logged-in" : "logged-out";
                userElement.innerHTML = `<span class="${statusClass}">●</span>${user.username}`;
                userElement.addEventListener("click", () => onUserClick(user.user_id));
                userContainer.appendChild(userElement);

                addedUserIds.add(user.user_id);
            }
        });
};