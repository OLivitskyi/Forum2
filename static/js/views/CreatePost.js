import AbstractView from "./AbstractView.js";
export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("CreatePost");
  }
  // If server side renders HTML, can use fetch API
  async getHtml() {
    return `
        <form class="form" id="create-post-view">
 <body>
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
          <a href="#" id="homepage">
            <span class="material-icons-sharp">grid_view</span>
            <h3>Dashboard</h3>
          </a>
          <a href="messages">
            <span class="material-icons-sharp">mail_outline</span>
            <h3>Messages</h3>
            <span class="message-count">26</span>
          </a>
          <a href="create-post" id="create-post" data-link>
            <span class="material-icons-sharp">add</span>
            <h3>Create Post</h3>
          </a>
          <a href="logout">
            <span class="material-icons-sharp" type= "submit">logout</span>
            <h3>Logout</h3>
          </a>
        </div>
      </aside>
      <!---- END OF ASIDE ---->
      <main>
        <h1>Create a post1</h1>
        <div class= "date">
        </div>
        <div class= "insights">
        </div>
      </main>
    </div>
  </body>
            </form>
            `;
  }
}