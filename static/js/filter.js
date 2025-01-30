import { createPostElement } from "./posts.js";


export const filterCat = (page = 1) => {
    let selectedCategories = [];
    const categoryList = document.querySelector(".categories");
    const postsContainer = document.querySelector("#posts-container");
    const myPostsFilter = document.getElementById('my-posts-filter');
    const myLikesFilter = document.getElementById('my-likes-filter');
    const loadMoreButton = document.getElementById('load-more');

    const fetchPosts = async (url) => {
        try {
            const resp = await fetch(url);
            if (!resp.ok) {
                console.log("Error fetching posts from API");
                return
            }

            const posts = await resp.json();
            if (posts.length === 0) {
                postsContainer.textContent = "No Posts Available";
                return
            }

            posts.forEach((post) => {
                const postElement = createPostElement(post)
                postsContainer.appendChild(postElement)
            })

            if (posts.length < 10) {
                if (loadMoreButton) {
                    loadMoreButton.disabled = true
                    loadMoreButton.style.opacity = 0.7
                }
            }
        } catch (error) {
            console.log("Error fetching posts:", error);
        }
    }

    let currentPage = page

    loadMoreButton?.addEventListener('click', () => {
        currentPage++;
        let url = selectedCategories.length > 0
            ? `/api/posts/${currentPage}/categories=${selectedCategories.join("&")}`
            : `/api/posts/${currentPage}`;
        fetchPosts(url);
    });

    myLikesFilter?.addEventListener('click', async () => {
        try {
            const resp = await fetch('/api/posts/liked');
            if (!resp.ok) {
                console.log("Didn't get liked posts from API");
                return;
            }

            const posts = await resp.json();
            postsContainer.replaceChildren();

            if (posts.length === 0) {
                postsContainer.textContent = "No Liked Posts Available";
                return;
            } else {
                posts.forEach((post) => {
                    const postElement = createPostElement(post);
                    postsContainer.appendChild(postElement);
                });
            }
        } catch (error) {
            console.error("Error fetching user liked posts:", error);
        }
    });

    myPostsFilter?.addEventListener('click', async () => {
        try {
            const resp = await fetch('/api/posts/created');
            if (!resp.ok) {
                console.log("Didn't get user posts from API");
                return;
            }

            const posts = await resp.json();
            postsContainer.replaceChildren();

            if (posts.length === 0) {
                postsContainer.textContent = "No Posts Available";
                return;
            } else {
                posts.forEach((post) => {
                    const postElement = createPostElement(post);
                    postsContainer.appendChild(postElement);
                });
            }
        } catch (error) {
            console.error("Error fetching user posts:", error);
        }
    });

    categoryList.addEventListener("click", async (event) => {
        if (event.target.tagName === "LI") {
            let value = event.target.textContent.trim();

            if (selectedCategories.includes(value)) {
                selectedCategories = selectedCategories.filter(cat => cat !== value);
                event.target.classList.remove("Selected", "active");
            } else {
                selectedCategories.push(value);
                event.target.classList.add("Selected", "active");
            }

            console.log("Selected categories:", selectedCategories);

            let url = selectedCategories.length > 0 
                ? `/api/posts/${page}/categories=${selectedCategories.join("&")}` 
                : `/api/posts/${page}`;

            try {
                const resp = await fetch(url);

                if (!resp.ok) {
                    console.log("Didn't get posts from API");
                    return;
                }

                const posts = await resp.json();
                postsContainer.replaceChildren();
                console.log("Posts:", posts);
                if (posts.length === 0) {
                    postsContainer.textContent = "No Posts Available";
                    return;
                } else {
                    posts.forEach((post) => {
                        const postElement = createPostElement(post);
                        postsContainer.appendChild(postElement);
                    });
                }
                
            } catch (error) {
                console.error("Error fetching posts:", error);
            }
        }
    });
};