import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Registration");
    }
    // If server side renders HTML, can use fetch API 
    async getHtml() { 
        return `
        <form class="form" id="createAccount">
        <div class="container-login">
         <h1 class="form__title">Create Account</h1>
            <div class="form__message form__message--error"></div>
            <div class="form__input-group">
    <input type="text" id="signupUsername" name="signupUsername" class="form__input" required autofocus placeholder="Username">
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <input type="number" id="quantity" min="18" class="form__input" name="age" autofocus placeholder="Age (must be 18 or older)">
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <select class="form__input" name="gender" autofocus placeholder="Gender">
        <option value="">Gender</option> <option value="male">Male</option>
        <option value="female">Female</option>
        <option value="other">Other</option>
    </select>
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <input type="text" class="form__input" name="firstname" autofocus placeholder="First name">
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <input type="text" class="form__input" name="lastname" autofocus placeholder="Last name">
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <input type="email" class="form__input" name="email" required autofocus placeholder="Email Address">
    <div class="form__input-error-message"></div>
</div>
<div class="form__input-group">
    <input type="password" class="form__input" name="signupPassword" id="pwd" required autofocus placeholder="Password">
    <div class="form__input-error-message"></div>
</div>
<button class="form__button" type="submit">Continue</button>
<p class="form__text">
    <a class="form__link" href="./" id="linkLogin">Already have an account? Sign in</a>
</p>

        </div>
           
</form>
        `;
    }
}

// Trying to add error message for @ character 



console.log("test registration.js")

document.getElementById("submit")
console.log(document.getElementById("submit"))