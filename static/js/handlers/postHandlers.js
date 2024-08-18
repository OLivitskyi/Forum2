import { navigateTo, navigateToPostDetails } from "../routeUtils.js";
import { createPost, getPosts, getPost } from "../api.js";
import { showError, clearError } from "../errorHandler.js";
import { sendPost } from "../websocket.js";
import { loadAndRenderComments } from "./commentHandlers.js";

export const handleCreatePostFormSubmit = () => {
    const form = document.getElementById("create-post-form");
    if (form) {
      console.log("Adding submit event listener to form");
      form.addEventListener("submit", handleCreatePostSubmit);
    } else {
      console.error("Form not found in handleCreatePostFormSubmit");
    }
  };
  
  async function handleCreatePostSubmit(e) {
    e.preventDefault();
    console.log("Handling post submit");
    clearError();
  
    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;
    const category = document.getElementById("category-select").value;
    const categoryError = document.getElementById("category-error");
  
    if (!title || !content || !category) {
      showError("Title, content, and category are required");
      if (!category) {
        categoryError.style.display = "block"; 
      } else {
        categoryError.style.display = "none"; 
      }
      return;
    }
  
    categoryError.style.display = "none";
  
    const body = { title, content, category_ids: [category] };
  
    try {
      console.log("Sending post data to server:", body);
      const response = await createPost(body);
  
      if (response.ok) {
        clearError();
        const post = await response.json();
        console.log("Post created successfully:", post);
        navigateTo("/homepage");
        sendPost(post);
      } else {
        const errorText = await response.text();
        console.error("Failed to create post:", errorText);
        showError(errorText || "Failed to create post");
      }
    } catch (error) {
      console.error("An error occurred while creating the post:", error);
      showError("An error occurred. Please try again.");
    }
  }
  
export const loadAndRenderSinglePost = async (postId) => {
    try {
        console.log("Loading post details for post ID:", postId);
        const postContainer = document.getElementById("single-post-container");
        if (!postContainer) {
            console.error("Post container not found");
            return;
        }
        const post = await getPost(postId);
        if (!post) {
            console.error("Post data is missing");
            return;
        }
        postContainer.innerHTML = `
            <div>
                <div class="post-container-title">${post.subject}</div>
                <div class="post-container-content">${post.content}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            </div>
        `;
        await loadAndRenderComments(postId);
    } catch (error) {
        console.error("Failed to load post:", error);
    }
};
export const loadAndRenderPosts = async () => {
    try {
        const postsContainer = document.getElementById("posts-container");
        if (!postsContainer) return;
        const posts = await getPosts();
        postsContainer.innerHTML = "";
        posts.forEach((post) => {
            const categories = post.categories
                .map((category) => `<span class="category">${category.name}</span>`)
                .join(", ");
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.id = `post-${post.id}`;
            postElement.innerHTML = `
                <h3>${post.subject}</h3>
                <p>${post.content}</p>
                <div class="post-categories">Categories: ${categories}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            `;
            postElement.addEventListener("click", () => {
                navigateToPostDetails(post.id);
            });
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Failed to load posts:", error);
    }
};
