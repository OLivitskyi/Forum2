export const getLayoutHtml = (content) => {
    return `
        <div class="container">
            <aside>
                <div class="top">
                    <div class="logo">
                        <img src="./images/logo.png" />
                        <h2>FOR<span class="danger">UM</span></h2>
                    </div>
                    <div class="close" id="close-btn">
                        <span class="material-icons-sharp">close</span>
                    </div>
                </div>
                <div class="sidebar">
                    <a href="homepage" data-link>
                        <span class="material-icons-sharp">grid_view</span>
                        <h3>Dashboard</h3>
                    </a>
                    <a href="messages" data-link>
                        <span class="material-icons-sharp">mail_outline</span>
                        <h3>Messages</h3>
                        <span class="message-count">27</span>
                    </a>
                    <a href="create-post" data-link>
                        <span class="material-icons-sharp">add</span>
                        <h3>Create Post</h3>
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
    `;
};
