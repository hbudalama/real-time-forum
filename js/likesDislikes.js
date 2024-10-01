export function initializeLikeDislikeButtons() {
    document.querySelectorAll('.like-button').forEach(button => {
        button.addEventListener('click', (event) => {
            console.log("you liked!!", event);
            event.preventDefault();
            const postId = button.getAttribute('data-id');

            fetch(`/api/posts/${postId}/like`, {
                method: 'POST',
                credentials: 'include'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    button.innerHTML = `<i class="fa fa-thumbs-up icon liked"></i>${data.likes}`;
                    const dislikeButton = button.closest('.post-row')?.querySelector('.dislike-button');
                    if (dislikeButton) {
                        dislikeButton.innerHTML = `<i class="fa fa-thumbs-down icon"></i>${data.dislikes}`;
                    }
                    
                } else {
                    console.error('Error:', data.message);
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });

    document.querySelectorAll('.dislike-button').forEach(button => {
        button.addEventListener('click', (event) => {
            event.preventDefault();
            const postId = button.getAttribute('data-id');

            fetch(`/api/posts/${postId}/dislike`, {
                method: 'POST',
                credentials: 'include'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    button.innerHTML = `<i class="fa fa-thumbs-down icon disliked"></i>${data.dislikes}`;
                    const likeButton = button.closest('.post-row')?.querySelector('.like-button');
                    if (likeButton) {
                        likeButton.innerHTML = `<i class="fa fa-thumbs-up icon"></i>${data.likes}`;
                    }
                    
                } else {
                    console.error('Error:', data.message);
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });
}
