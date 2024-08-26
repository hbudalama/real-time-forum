function initializeComments() {
    const mainContent = document.getElementById('main-content');
    const dialog = document.getElementById('comment-dialog');
    const dialogOverlay = document.getElementById('dialog-overlay');
    const commentList = document.getElementById('comment-list');
    const postTitleElement = document.getElementById('post-title');
    const postContentElement = document.getElementById('post-content');
    const commentTextarea = document.querySelector('#comment-dialog textarea[name="comment-area"]');
    const addCommentButton = document.getElementById('add-comment-button');
    const commentsCountElement = document.getElementById('comments-count'); // Assuming this element exists

    const renderComments = (comments) => {
        if (!Array.isArray(comments) || comments.length === 0) {
            commentList.innerHTML = '<p>No comments yet. Be the first to comment!</p>';
        } else {
            const commentItems = comments.map(comment => `
                <li>
                    <p><strong>${comment.Username}:</strong> ${comment.Content}</p>
                    <p>${new Date(comment.CreatedDate).toLocaleString()}</p>
                    <div class="post-row">
                        <div class="activity-icons">
                            <div class="comment-like-button" data-id="${comment.ID}">
                                <i class="fa fa-thumbs-up icon"></i>${comment.Likes}
                            </div>
                            <div class="comment-dislike-button" data-id="${comment.ID}">
                                <i class="fa fa-thumbs-down icon"></i>${comment.Dislikes}
                            </div>
                        </div>
                    </div>
                </li>
            `).join('');
            commentList.innerHTML = commentItems;
        }
    
        // Update comments count
        if (commentsCountElement) {
            commentsCountElement.textContent = `Comments (${comments.length})`;
        }
    
        // Initialize the like and dislike buttons for comments
        initializeCommentLikeDislikeButtons();
    };

    const initializeCommentLikeDislikeButtons = () => {
        document.querySelectorAll('.comment-like-button').forEach(button => {
            button.addEventListener('click', (event) => {
                event.preventDefault();
                const commentId = button.getAttribute('data-id');
    
                fetch(`/api/comments/${commentId}/like`, {
                    method: 'POST',
                    credentials: 'include'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        button.innerHTML = `<i class="fa fa-thumbs-up icon"></i>${data.likes}`;
                        const dislikeButton = button.closest('.post-row').querySelector('.comment-dislike-button');
                        if (dislikeButton) {
                            dislikeButton.innerHTML = `<i class="fa fa-thumbs-down icon"></i>${data.dislikes}`;
                        }
                    } else {
                        console.error('Error:', data.reason);
                    }
                })
                .catch(error => console.error('Error:', error));
            });
        });
    
        document.querySelectorAll('.comment-dislike-button').forEach(button => {
            button.addEventListener('click', (event) => {
                event.preventDefault();
                const commentId = button.getAttribute('data-id');
    
                fetch(`/api/comments/${commentId}/dislike`, {
                    method: 'POST',
                    credentials: 'include'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        button.innerHTML = `<i class="fa fa-thumbs-down icon"></i>${data.dislikes}`;
                        const likeButton = button.closest('.post-row').querySelector('.comment-like-button');
                        if (likeButton) {
                            likeButton.innerHTML = `<i class="fa fa-thumbs-up icon"></i>${data.likes}`;
                        }
                    } else {
                        console.error('Error:', data.reason);
                    }
                })
                .catch(error => console.error('Error:', error));
            });
        });
    };

    const addComment = () => {
        const comment = commentTextarea.value.trim();
        if (comment === '') {
            alert('Comment cannot be empty');
            return;
        }
        const postId = dialog.dataset.postId; // Assuming postId is set as a data attribute on the dialog
        
        fetch(`/api/posts/${postId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: new URLSearchParams({
                comment: comment
            }),
            credentials: 'include'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Add the new comment to the list
                renderComments(data.comments);
                commentTextarea.value = ''; // Clear the textarea
            } else {
                console.error('Error:', data.message);
            }
        })
        .catch(error => console.error('Error:', error));
    };

    addCommentButton.addEventListener('click', addComment);

    mainContent.addEventListener('click', (event) => {
        const postLink = event.target.closest('.post-title-link');
        const commentIcon = event.target.closest('.comment-icon');

        if (postLink || commentIcon) {
            const postId = postLink ? postLink.dataset.id : commentIcon.dataset.id;

            fetch(`/api/posts/${postId}/comments`)
                .then(response => response.json())
                .then(data => {
                    postTitleElement.textContent = data.post.Title;
                    postContentElement.textContent = data.post.Content;
                    renderComments(data.comments);

                    // Set postId as a data attribute on the dialog
                    dialog.dataset.postId = postId;

                    dialog.classList.add('show');
                    dialogOverlay.classList.add('show');
                })
                .catch(error => console.error('Error fetching post details and comments:', error));
        }
    });

    dialogOverlay.addEventListener('click', () => {
        dialog.classList.remove('show');
        dialogOverlay.classList.remove('show');
    });
}