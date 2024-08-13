export const getLayoutHtml = (content) => {
    return `
        <div class="container">
            <aside>
                <div class="top">
                    <div class="logo">
                        <span class="material-icons-sharp logo-icon">app_registration</span>
                        <h2>FOR<span class="danger">UM</span></h2>
                    </div>
                    <div class="close" id="close-btn">
                        <span class="material-icons-sharp">close</span>
                    </div>
                </div>
                <div class="sidebar">
                    <a href="homepage" data-link id="homepage">
                        <span class="material-icons-sharp">grid_view</span>
                        <h3>Dashboard</h3>
                    </a>
                    <a href="messages" data-link id="messages">
                        <span class="material-icons-sharp">mail_outline</span>
                        <h3>Messages</h3>
                        <span class="message-count">0</span>
                    </a>
                    <a href="create-post" data-link id="create-post">
                        <span class="material-icons-sharp">add</span>
                        <h3>Create Post</h3>
                    </a>
                    <a href="create-category" data-link id="create-category">
                        <span class="material-icons-sharp">category</span>
                        <h3>Create Category</h3>
                    </a>
                    <a href="#" id="logout">
                        <span class="material-icons-sharp">logout</span>
                        <h3>Logout</h3>
                    </a>
                </div>
            </aside>
            <main>
                ${content}
            </main>
        </div>
        <div id="popup-notification">You have a new message!</div>
    `;
};
