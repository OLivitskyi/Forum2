import { handleLoginFormSubmit, handleLogout } from './handlers/authHandlers.js';
import { handleCreatePostFormSubmit, loadAndRenderPosts,loadAndRendeSinglePost } from './handlers/postHandlers.js';
import { handleCreateCategoryFormSubmit, loadCategories } from './handlers/categoryHandlers.js';
import { setupWebSocketHandlers } from './websocket.js';

export { 
    handleLoginFormSubmit, 
    handleLogout, 
    handleCreatePostFormSubmit, 
    loadAndRenderPosts, 
    handleCreateCategoryFormSubmit, 
    loadCategories, 
    setupWebSocketHandlers,
    loadAndRendeSinglePost
};
