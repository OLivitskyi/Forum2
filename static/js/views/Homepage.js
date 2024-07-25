import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Homepage");
  }
  // If server side renders HTML, can use fetch API
  async getHtml() {
    return `
        <form class="form" id="homepage">
            <body>
                <div class="container-sidebar">
                    <aside>
                        <div class="top">
                                <div class="logo">
                                    <span class="material-icons-sharp">save_as</span>
                                    <h2>FOR<span class="danger">UM</span></h2>
                                </div>
                                <div class="close" id="close-btn">
                                    <span class="material-icons-sharp">close</span>
                                </div>
                        </div>
                        <div class="sidebar">
                              <a href="#/dashboard">
                                <span class="material-icons-sharp">grid_view</span>
                                <h3>Dashboard</h3>
                              </a>
                              <a href="#" class="active">
                                <span class="material-icons-sharp">person_outline</span>
                                <h3>Profile</h3>
                              </a>
                              <a href="#">
                                <span class="material-icons-sharp">mail_outline</span>
                                <h3>Messages</h3>
                                <span class="message-count">26</span>
                              </a>
                              <a href="#">
                                <span class="material-icons-sharp">add</span>
                                <h3>Create Post</h3>
                              </a>
                              <a href="/logout" id="logout" data-link>
                                <span class="material-icons-sharp">logout</span>
                                <h3>Logout</h3>
                              </a>
                        </div>
                  </aside>
    
              </div>
              <main>
              <div class= "top-bar">
                  <h1>Welcome to your Dashboard, User!</h1>
              </div>
              </main>  
  </body>
            </form>
            `;
  }
}
