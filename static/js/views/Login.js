import AbstractView from "./AbstractView.js";
export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Login");
  }
  async getHtml() {
    return `
      <form class="form" id="login">
      <div class="container-login">
                <h1 class="form__title">Login</h1>
          <div id="error-message" class="form__message form__message--error"></div>
          <div class="form__input-group">
              <input type="text" class="form__input" name="username" autofocus placeholder="Username or email">
              <div class="form__input-error-message"></div>
          </div>
          <div class="form__input-group">
              <input type="password" class="form__input" name="password" autofocus placeholder="Password">
              <div class="form__input-error-message"></div>
          </div>
          <button class="form__button" type="submit">Continue</button>
          <p class="form__text">
              <a href="#" class="form__link">Forgot your password?</a>
          </p>
          <p class="form__text">
              <a class="form__link" href="./registration" id="linkCreateAccount" data-link>Don't have an account? Create account</a>
          </p>
                </div>
      </form>
      `;
  }
}
