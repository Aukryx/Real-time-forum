import { createWelcomePage, removeWelcomePage } from "./welcome.js";

let logoutButton = document.getElementById('logout')

if (logoutButton) {
    logoutButton.addEventListener('click', () => {
        removeWelcomePage()
        createWelcomePage()
    })
}