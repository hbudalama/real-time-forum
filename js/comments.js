function initializeComments() {
    const mainContent = document.getElementById('main-content');
    const dialog = document.getElementById('comment-dialog');
    const dialogOverlay = document.getElementById('dialog-overlay');
    const commentList = document.getElementById('comment-list');
    const postTitleElement = document.getElementById('post-title');
    const postContentElement = document.getElementById('post-content');

    const renderComments = (comments) => {
        if (!Array.isArray(comments) || comments.length === 0) {
            commentList.innerHTML = '<p>No comments yet. Be the first to comment!</p>';
            return;
        }
    
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
    
        // Initialize the like and dislike buttons for comments
        initializeCommentLikeDislikeButtons();
    };
    

    const initializeCommentLikeDislikeButtons = () => {
        document.querySelectorAll('.comment-like-button').forEach(button => {
            button.addEventListener('click', (event) => {
                event.preventDefault(); // Prevent the default action
                const commentId = button.getAttribute('data-id');
    
                fetch(`/api/comments/${commentId}/like`, {
                    method: 'POST',
                    credentials: 'include'
                })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            console.error('Server error:', text);
                            throw new Error(`Server error: ${response.status} ${response.statusText}`);
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        // Update both like and dislike counts
                        button.innerHTML = `<i class="fa fa-thumbs-up icon"></i>${data.likes}`;
                        const dislikeButton = button.closest('.post-row').querySelector('.comment-dislike-button');
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
    
        document.querySelectorAll('.comment-dislike-button').forEach(button => {
            button.addEventListener('click', (event) => {
                event.preventDefault(); // Prevent the default action
                const commentId = button.getAttribute('data-id');
    
                fetch(`/api/comments/${commentId}/dislike`, {
                    method: 'POST',
                    credentials: 'include'
                })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            console.error('Server error:', text);
                            throw new Error(`Server error: ${response.status} ${response.statusText}`);
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        // Update both dislike and like counts
                        button.innerHTML = `<i class="fa fa-thumbs-down icon"></i>${data.dislikes}`;
                        const likeButton = button.closest('.post-row').querySelector('.comment-like-button');
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
    };    

    mainContent.addEventListener('click', (event) => {
        const postLink = event.target.closest('.post-title-link');
        const commentIcon = event.target.closest('.comment-icon');

        if (postLink || commentIcon) {
            const postId = postLink ? postLink.dataset.id : commentIcon.dataset.id;

            fetch(`/api/posts/${postId}/comments`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`Failed to fetch post details and comments: ${response.status} ${response.statusText}`);
                    }
                    return response.json();
                })
                .then(data => {
                    // Update the post title and content in the dialog
                    postTitleElement.textContent = data.post.Title;
                    postContentElement.textContent = data.post.Content;

                    // Render the comments
                    renderComments(data.comments);

                    // Show the dialog
                    dialog.classList.add('show');
                    dialogOverlay.classList.add('show');
                })
                .catch(error => {
                    console.error('Error fetching post details and comments:', error);
                });
        }
    });

    dialogOverlay.addEventListener('click', () => {
        dialog.classList.remove('show');
        dialogOverlay.classList.remove('show');
    });
}
