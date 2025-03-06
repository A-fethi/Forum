import { showNotification } from "./components/notifications.js";
import { createPostElement } from "./posts.js";

export let post = {};

export function checkPost() {
  const isLoggedIn = document.cookie.includes("IsLoggedIn=true");
  const createPostButton = document.getElementById("create-post-button");

  if (createPostButton) {
    createPostButton.addEventListener("click", async (event) => {
      event.preventDefault();
      const titleInput = document.querySelector('input[name="title"]');
      const contentInput = document.querySelector('textarea[name="content"]');
      const selectedCategories = Array.from(
        document.querySelectorAll('input[name="category"]:checked')
      ).map((checkbox) => checkbox.value);
      if (titleInput && contentInput) {
        let title = titleInput.value.trim();
        let content = contentInput.value.trim();
        const data = {
          title: title,
          content: content,
          categories: selectedCategories,
        };

        try {
          const resp = await fetch("/api/posts/add", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
            credentials: "include",
          });

          if (resp.status === 201) {
            const responseData = await resp.json();
            titleInput.value = "";
            contentInput.value = "";
            document
              .querySelectorAll('input[name="category"]:checked')
              .forEach((checkbox) => (checkbox.checked = false));
            const postsElement = document.getElementById("posts-container");
            postsElement.prepend(createPostElement(responseData));
          } else {
            const responseData = await resp.json();
            console.error("Failed to create post:", resp.statusText);
            showNotification(responseData.message, "error");
          }
        } catch (error) {
          console.error("Error occurred while creating post:", error);
          showNotification(
            "An error occurred, Please try again later",
            "error"
          );
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
