// display and add functions to login & signup page
import { globals, eventListeners } from './globals.js';
import { renderHome } from './render_home.js';
export { displayLogin };

// LOGIN ----------------------------------------------------

window.addEventListener('DOMContentLoaded', async function () {

    console.log(globals);
    console.log(eventListeners);

    const isLoggedIn = await checkLogin();
    if (isLoggedIn) {
        renderHome();
    } else {
        displayLogin();
    }
});

function displayLogin() {
    if (!eventListeners.addEventToLoginForm) {
        addEventToLoginForm();
        eventListeners.addEventToLoginForm = true;
    }
    if (!eventListeners.addEventSignupLink) {
        addEventSignupLink();
        eventListeners.addEventSignupLink = true;
    }

    globals.currentPage = "login";
    document.getElementById("login").style.display = 'block';
    document.getElementById("home").style.display = 'none';
    document.getElementById("signup").style.display = 'none';
}


function addEventToLoginForm() {
    document.getElementById("login_form").addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(event.target);

        return fetch('/login', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.loginStatus === "success") {
                    globals.LoggedInUser = data.user;
                    globals.currentFilter = "all"
                    location.reload();
                    renderHome();
                } else if (data.loginStatus === "fail") {
                    alert(data.message);
                }
            })
            .catch(error => {
                console.error('Error logging in:', error);
            });
    });
}

function addEventSignupLink() {
    document.getElementById('signup_link').addEventListener('click', function () {
        if (!eventListeners.addEventToSignUpForm) {
            addEventToSignupForm();
            eventListeners.addEventToSignUpForm = true;
        }
        if (!eventListeners.addEventLoginLink) {
            addEventLoginLink();
            eventListeners.addEventLoginLink = true;
        }

        globals.currentPage = "signup";
        document.getElementById("login").style.display = 'none';
        document.getElementById("home").style.display = 'none';
        document.getElementById("signup").style.display = 'block';
    });
}

async function checkLogin() {
    try {
        const response = await fetch('/check-login')
        const data = await response.json();
        globals.LoggedInUser = data.user;
        return true;
    } catch (error) {
        console.error('Error checking login:', error);
        return false;
    }
}

// SIGNUP ----------------------------------------------------

function addEventLoginLink() {
    document.getElementById('login_link').addEventListener('click', function () {
        displayLogin();
    });
}

function addEventToSignupForm() {
    document.getElementById("signup_form").addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(event.target);

        fetch('/signup', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.status === "success") {
                    alert(data.message);
                    location.reload();
                } else if (data.status === "error") {
                    alert(data.message);
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
    });
}
