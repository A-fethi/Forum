import { showNotification } from "./components/notifications.js";
import { filterCat } from "./filter.js";
import { createPostElement } from "./posts.js";
import { createCommentElement } from "./comments.js";
import { logout } from "./auth.js";
import { checkPost } from "./addposts.js";
export const Home = async () => {
    console.log("Initializing Home function...");
    const loadMoreButton = document.getElementById('load-more');
    console.log(loadMoreButton)
    if (loadMoreButton) {
        loadMoreButton.style.display = 'none';
    }


    logout();

    checkPost();
    console.log("Checked post creation functionality.");
    const postsElement = document.getElementById("posts-container");
    console.log("Posts container element:", postsElement);

    try {
        console.log("Fetching posts from /api/posts...");
        const resp = await fetch("/api/posts");

        if (!resp.ok) {
            const res = await resp.json()
            console.error("Failed to fetch posts, response not ok.");
            // postsElement.textContent = "No Posts Found";
            showNotification(res.message, "error");
            return;
        }

        const posts = await resp.json();
        console.log("Fetched posts:", posts);

        if (!posts || posts.length === 0) {
            loadMoreButton.style.display = 'none';
            console.warn("No posts available.");
            showNotification("No posts found", "error");
            return;
        }
        if (loadMoreButton) {
            if (posts.length < 10) {
                loadMoreButton.style.display = 'none';
            } else {
                loadMoreButton.style.display = 'block';
            }
        }
        postsElement.replaceChildren();
        posts.forEach((post) => {
            console.log("Processing post:", post);
            const postElement = createPostElement(post);
            postsElement.appendChild(postElement);
        });
        
        console.log("Comment functionality initialized.");

        filterCat();
        console.log("Category filtering initialized.");

        console.log("Interaction listeners attached.");

    } catch (error) {
        console.error("Error: ", error);
        // // postsElement.textContent = "No Posts Available";
        // showNotification("an Error Occured, Try Again", "error");
    }
};
