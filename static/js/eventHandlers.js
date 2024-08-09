import { handleLoginFormSubmit, handleLogout } from './handlers/authHandlers.js';
import { handleCreatePostFormSubmit, loadAndRenderPosts, loadAndRenderSinglePost} from './handlers/postHandlers.js';
import { loadAndRenderComments } from './handlers/commentHandlers.js';
import { handleCreateCategoryFormSubmit, loadCategories } from './handlers/categoryHandlers.js';
import { setupWebSocketHandlers } from './websocket.js';
import { navigateToPostDetails } from './routeUtils.js';


export { 
    handleLoginFormSubmit, 
    handleLogout, 
    handleCreatePostFormSubmit, 
    loadAndRenderPosts, 
    handleCreateCategoryFormSubmit, 
    loadCategories, 
    setupWebSocketHandlers,
    loadAndRenderSinglePost,
    loadAndRenderComments,
    navigateToPostDetails
};