import { openAuthModal } from "./auth.js";
import { logout } from "./auth.js";
import { Home } from "./Home.js";
import { checkPost } from "./addposts.js";
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

// export function getCookie(name) {
//     console.log(`Retrieving cookie: ${name}`);
//     const value = `; ${document.cookie}`;
//     const parts = value.split(`; ${name}=`);
//     if (parts.length === 2) {
//         console.log(`Cookie ${name} found.`);
//         return parts.pop().split(';').shift();
//     } else {
//         console.warn(`Cookie ${name} not found.`);
//         return null;
//     }
// }

console.log("Initializing home and session...");
Home();
Session();
logout();
SessionCheck();
console.log("Session check initialized.");
