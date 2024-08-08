import { navigateTo, navigateToPostDetails } from '../router.js';
import { createPost, getPosts } from '../api.js';
import { showError, clearError } from '../errorHandler.js';
import { sendPost } from '../websocket.js';

export const handleCreatePostFormSubmit = (clearError, showError) => {
    const form = document.getElementById("create-post-form");
    if (form) {
        form.removeEventListener("submit", handleCreatePostSubmit);
        form.addEventListener("submit", handleCreatePostSubmit, { once: true });
    }
};

async function handleCreatePostSubmit(e) {
    e.preventDefault();
    clearError();
    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;
    const categorySelect = document.getElementById("category-select");
    const category = categorySelect.value;
    if (!title || !content || !category) {
        showError("Title, content, and category are required");
        return;
    }
    const body = {
        title,
        content,
        category_ids: [category]
    };
    try {
        const response = await createPost(body);
        if (response.ok) {
            clearError();
            navigateTo("/homepage");

            const post = await response.json();
            sendPost(post); // Відправка посту через вебсокет
        } else {
            const errorText = await response.text();
            showError(errorText || "Failed to create post");
        }
    } catch (error) {
        showError("An error occurred. Please try again.");
    }
};

export const loadAndRenderPosts = async () => {
    try {
        const postsContainer = document.getElementById("posts-container");
        if (!postsContainer) return;
        const posts = await getPosts();
        postsContainer.innerHTML = "";
        posts.forEach(post => {
            const categories = post.categories.map(category => `<span class="category">${category.name}</span>`).join(', ');
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <h3>${post.subject}</h3>
                <p>${post.content}</p>
                <div class="post-categories">Categories: ${categories}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            `;
            postElement.onclick = () => navigateToPostDetails(post.id);
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Failed to load posts:", error);
    }
};
