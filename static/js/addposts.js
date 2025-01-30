import { openAuthModal } from "./auth.js";
import { showNotification } from "./components/notifications.js";
import { Home } from "./Home.js";
import { createPostElement } from "./posts.js";

export let post = {};

export function checkPost() {
    console.log("Checking post creation setup...");
    const isLoggedIn = document.cookie.includes("IsLoggedIn=true");
    console.log("Is user logged in?", isLoggedIn);

    const createPostButton = document.getElementById("create-post-button");

    if (createPostButton) {
        console.log("Submit button found, adding event listener.");

        createPostButton.addEventListener('click', async (event) => {
            event.preventDefault();
            console.log("Submit button clicked.");

            const titleInput = document.querySelector('input[name="title"]');
            const contentInput = document.querySelector('textarea[name="content"]');
            const selectedCategories = Array.from(
                document.querySelectorAll('input[name="category"]:checked')
            ).map((checkbox) => checkbox.value);


            console.log("Collected form data:", {
                title: titleInput?.value,
                content: contentInput?.value,
                categories: selectedCategories
            });


            if (titleInput && contentInput) {
                let title = titleInput.value.trim();
                let content = contentInput.value.trim();
                const data = {
                    title: title,
                    content: content,
                    categories: selectedCategories,
                };

                console.log("Sending post data:", data);

                try {
                    const resp = await fetch('/api/posts/', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(data),
                        credentials: "include",
                    });

                    if (resp.status === 201) {
                        const responseData = await resp.json();
                        console.log('Post created successfully:', responseData);
                        titleInput.value = '';
                        contentInput.value = '';
                        document.querySelectorAll('input[name="category"]:checked').forEach(checkbox => checkbox.checked = false);
                        const postsElement = document.getElementById("posts-container");
                        postsElement.prepend(createPostElement(responseData));
                        
                    } else {
                        const responseData = await resp.json();
                        console.error('Failed to create post:', resp.statusText);
                        showNotification(responseData.message, "error");
                        // openAuthModal();
                    }
                } catch (error) {
                    console.error("Error occurred while creating post:", error);
                    showNotification("An error occurred, Please try again later", "error");
                }
            } else {
                console.error("Title or Content inputs not found.");
                showNotification("Error: Title/Content cannot be empty", "error");
            }
        });
    } else {
        console.error("Submit button not found.");
    }
}
