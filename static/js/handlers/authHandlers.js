import { navigateTo } from '../routeUtils.js';
import { sendRequest } from '../api.js';
import { connectWebSocket, setupWebSocketHandlers } from '../websocket.js';
import { showError, clearError } from '../errorHandler.js';
import { connectAfterLogin } from '../auth.js';

export const handleLoginFormSubmit = () => {
    const loginForm = document.getElementById("login");
    if (loginForm) {
        loginForm.removeEventListener("submit", handleLogin);
        loginForm.addEventListener("submit", handleLogin, { once: true });
        console.log("Login form handler attached");
    }
};

async function handleLogin(e) {
    e.preventDefault();
    clearError();
    const formData = new FormData(e.target);
    try {
        const response = await fetch("/api/login", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            const token = await response.text();
            console.log("Login successful, token received:", token);
            connectAfterLogin(token);
            clearError();
            navigateTo("/homepage");
        } else {
            const errorText = await response.text();
            console.error("Login failed, server response:", errorText);
            showError(errorText || "Login failed. Please try again.");
        }
    } catch (error) {
        console.error("An error occurred during login:", error);
        showError("An error occurred. Please try again.");
    }
};

export const handleLogout = () => {
    const logoutButton = document.getElementById("logout");
    if (logoutButton) {
        logoutButton.removeEventListener("click", handleLogoutClick);
        logoutButton.addEventListener("click", handleLogoutClick, { once: true });
    }
};

async function handleLogoutClick(e) {
    e.preventDefault();
    await logout();
};

export const logout = async () => {
    const response = await sendRequest("/logout", "GET");
    if (response.ok) {
        clearCookies();
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};

function clearCookies() {
    const cookies = document.cookie.split(";");
    for (const cookie of cookies) {
        const eqPos = cookie.indexOf("=");
        const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/;";
        console.log(`Cleared cookie: ${name}`);
    }
}
