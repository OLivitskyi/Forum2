import Homepage from "./views/Homepage.js";
import Login from "./views/Login.js";
import Registration from "./views/Registration.js";
import CreatePost from "./views/CreatePost.js";
import Messages from "./views/Messages.js";
import CreatePostCategory from "./views/CreatePostCategory.js";
import PostDetails from "./views/Post.js";

export const routes = [
    { path: "/", view: Login },
    { path: "/registration", view: Registration },
    { path: "/homepage", view: Homepage, protected: true },
    { path: "/logout", view: Login },
    { path: "/create-post", view: CreatePost, protected: true },
    { path: "/messages", view: Messages, protected: true },
    { path: "/create-category", view: CreatePostCategory, protected: true },
    { path: "/post/:id", view: PostDetails }
];
