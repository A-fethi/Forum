import { openAuthModal } from "./auth.js";
import { logout } from "./auth.js";
import { Home } from "./Home.js";
import { SessionCheck } from "./components/sessionChecker.js";

export const Session = () => {
    console.log("Initializing session...");
    const authbtn = document.getElementById("login-register");
    if (authbtn) {
        console.log("Login/Register button found, adding event listener.");
        authbtn.addEventListener("click", (e) => {
            console.log("Login/Register button clicked.");
            openAuthModal();
        });
    } 
};

console.log("Initializing home and session...");
Home();
Session();
logout();
SessionCheck();
console.log("Session check initialized.");
