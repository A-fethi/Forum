let likeButton = document.getElementById('like-btn');
let likeCount = document.getElementById('likes-count');
let count = 0;

likeButton.addEventListener('click', () => {
    count++
    likeCount.innerText = count;
});