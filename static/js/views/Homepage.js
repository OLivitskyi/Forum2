import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Homepage");
    }

    async getHtml() {
        return `
            <form class="form" id="homepage-view">
                <body>
                    <div class="container">
                        <aside>
                            <div class="top">
                                <div class="logo">
                                    <span class="material-icons-sharp forum-logo">groups</span>
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
                                <a href="create-post" id="create-post" data-link>
                                    <span class="material-icons-sharp">add</span>
                                    <h3>Create Post</h3>
                                </a>
                                <a href="logout" data-link>
                                    <span class="material-icons-sharp" type="submit">logout</span>
                                    <h3>Logout</h3>
                                </a>
                            </div>
                        </aside>
                        <!---- END OF ASIDE ---->
                        <main>
                            <h1>Have something to share?</h1>
                            <div class="create-post-container">
                                <h2>What's on your mind?</h2>
                                <div class="container-post">
                                    <input type="text" name="post-title" class="post-title" placeholder="Title of your post">
                                    <input type="text" name="post-content" class="post-subject" placeholder="Write a post">
                                </div>
                                <button class="pill" type="button">General</button>
                                <button class="pill pill--selected" type="button">Travel</button>
                                <button class="pill" type="button">Hobbies</button>
                                <button class="pill" type="button">Gaming</button>
                                <button class="pill pill-submit" type="button">POST</button>
                            </div>
                        </main>
                    </div>
                </body>
            </form>
        `;
    }
}
