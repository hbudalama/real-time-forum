function initializeComments() {
    console.log("initializing comments")
    const mainContent = document.getElementById('main-content');
    if (!mainContent) {
        console.error("Error: 'main-content' element not found.");
        return;
    }

    const dialog = document.getElementById('comment-dialog');
    if (!dialog) {
        console.error("Error: 'comment-dialog' element not found.");
        return;
    }

    const dialogOverlay = document.getElementById('dialog-overlay');
    if (!dialogOverlay) {
        console.error("Error: 'dialog-overlay' element not found.");
        return;
    }

    const commentList = document.getElementById('comment-list');
    if (!commentList) {
        console.error("Error: 'comment-list' element not found.");
        return;
    }

    const commentTextarea = document.getElementById('comment-area');
    if (!commentTextarea) {
        console.error("Error: 'comment-area' textarea not found.");
        return;
    }

    const addCommentButton = document.getElementById('add-comment-button');
    if (!addCommentButton) {
        console.error("Error: 'add-comment-button' element not found.");
        return;
    }

    // Remove any existing event listener before adding a new one
    addCommentButton.removeEventListener('click', addComment);
    addCommentButton.addEventListener('click', addComment);
    console.log("Adding comment listener")


    mainContent.removeEventListener('click', openComments);
    mainContent.addEventListener('click', openComments);

    dialogOverlay.addEventListener('click', () => {
        dialog.classList.remove('show');
        dialogOverlay.classList.remove('show');
    });
}

/**
 * 
 * @param {MouseEvent} event 
 * @returns 
 */
function openComments(event) {
    const postTitleElement = document.getElementById('post-title');
    if (!postTitleElement) {
        console.error("Error: 'post-title' element not found.");
        return;
    }

    const postContentElement = document.getElementById('post-content');
    if (!postContentElement) {
        console.error("Error: 'post-content' element not found.");
        return;
    }

    const commentList = document.getElementById('comment-list');
    if (!commentList) {
        console.error("Error: 'comment-list' element not found.");
        return;
    }

    const dialog = document.getElementById('comment-dialog');
    if (!dialog) {
        console.error("Error: 'comment-dialog' element not found.");
        return;
    }

    const dialogOverlay = document.getElementById('dialog-overlay');
    if (!dialogOverlay) {
        console.error("Error: 'dialog-overlay' element not found.");
        return;
    }

    if (!event.target) {
        console.error("event target is null");
        return;
    }
    
    // @ts-ignore
    const target = event.target;

    const postLink = target.closest('.post-title-link');
    const commentIcon = target.closest('.comment-icon');

    if (!postLink && !commentIcon) {
        return;
    }
    
    // @ts-ignore
    const postId = postLink ? postLink.dataset.id : commentIcon.dataset.id;

    fetch(`/api/posts/${postId}/comments`)
        .then(response => response.json())
        .then(data => {
            postTitleElement.textContent = data.post.Title;
            postContentElement.textContent = data.post.Content;
            renderComments(data.comments, commentList);

            // Set postId as a data attribute on the dialog
            dialog.dataset.postId = postId;

            dialog.classList.add('show');
            dialogOverlay.classList.add('show');
        })
        .catch(error => console.error('Error fetching post details and comments:', error));
}

/**
 * 
 * @param {unknown} comments 
 * @param {HTMLElement} commentList
 */
function renderComments(comments, commentList) {
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
    initializeCommentLikeDislikeButtons();
};

function initializeCommentLikeDislikeButtons() {
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

/**
 * @returns 
 */
function addComment() {
    const dialog = document.getElementById('comment-dialog');
    const commentList = document.getElementById('comment-list');
    const commentTextarea = document.getElementById('comment-area')

    if (!dialog) {
        console.error("Error: 'comment-dialog' element not found.");
        return;
    }

    if (!commentList) {
        console.error("Error: 'comment-list' element not found.");
        return;
    }

    if (!commentTextarea) {
        console.error("Error: 'comment-area' element not found.");
        return;
    }

    // @ts-ignore
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
                const commentsCount = document.getElementById(`post-${postId}-comments-count`);
                if (commentsCount) {
                    commentsCount.innerText = (Number.parseInt(commentsCount.innerText) + 1).toString();
                }
                renderComments(data.comments, commentList);
                // @ts-ignore
                commentTextarea.value = ''; // Clear the textarea
            } else {
                console.error('Error:', data.message);
            }
        })
        .catch(error => console.error('Error:', error));
};