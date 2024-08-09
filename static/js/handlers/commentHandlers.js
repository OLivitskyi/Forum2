import { getComments, createComment } from '../api.js';
import { sendComment } from '../websocket.js';
import { setFormMessage } from '../formHandler.js';
import { showError, clearError } from '../errorHandler.js';

export const handleCreateCommentFormSubmit = (postId, clearError, showError) => {
    const form = document.getElementById("create-comment-form");
    if (form) {
        form.removeEventListener("submit", (e) => handleCreateCommentSubmit(e, postId));
        form.addEventListener("submit", (e) => handleCreateCommentSubmit(e, postId), { once: true });
    }
};

async function handleCreateCommentSubmit(e, postId) {
    e.preventDefault();
    clearError();
    const content = document.getElementById("content").value;
    const messageElement = document.getElementById("comment-message");
    if (!content) {
        showError("Content is required");
        return;
    }
    
    console.log("Submitting comment with post ID:", postId, "and content:", content);
    
    const body = {
        post_id: postId,
        content
    };
    
    try {
        const response = await createComment(body);
        if (response) {
            clearError();
            console.log("Comment created successfully.");

            // Показуємо повідомлення про успіх
            setFormMessage(messageElement, "success", "Comment added successfully");

            // Очищуємо поле введення
            document.getElementById("content").value = "";

            // Відправка коментаря через WebSocket
            sendComment(response); // Передаємо об'єкт коментаря з ID
        }
    } catch (error) {
        console.error("An error occurred while creating comment:", error);
        setFormMessage(messageElement, "error", "An error occurred. Please try again.");
    }
}


export const loadAndRenderComments = async (postId) => {
    try {
        const commentsContainer = document.getElementById("comments-container");
        if (!commentsContainer) return;
        const comments = await getComments(postId);
        commentsContainer.innerHTML = "";
        comments.forEach(comment => {
            const commentElement = document.createElement("div");
            commentElement.classList.add("comment");
            commentElement.innerHTML = `
                <h4>${comment.user.username}</h4>
                <p>${comment.content}</p>
                <div>
                    <span>Likes: ${comment.like_count}</span>
                    <span>Dislikes: ${comment.dislike_count}</span>
                </div>
            `;
            commentsContainer.appendChild(commentElement);
        });
    } catch (error) {
        console.error("Failed to load comments:", error);
    }
};