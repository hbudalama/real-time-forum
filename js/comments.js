function initializeComments() {
    const mainContent = document.getElementById('main-content');
    const dialog = document.getElementById('comment-dialog');
    const dialogOverlay = document.getElementById('dialog-overlay');
    const commentList = document.getElementById('comment-list');


    const renderComments = (comments) => {
        const commentItems = comments.map(comment => `
            <li>
                <p><strong>${comment.Username}:</strong> ${comment.Content}</p>
                <p>${new Date(comment.CreatedDate).toLocaleString()}</p>
                <div class="post-row">
                    <div class="activity-icons">
                        <div>
                            <a href="/api/comments/${comment.ID}/like">
                                <i class="fa fa-thumbs-up icon"></i>${comment.Likes}
                            </a>
                        </div>
                        <div>
                            <a href="/api/comments/${comment.ID}/dislike">
                                <i class="fa fa-thumbs-down icon"></i>${comment.Dislikes}
                            </a>
                        </div>
                    </div>
                </div>
            </li>
        `).join('');
        commentList.innerHTML = commentItems;
    };

    mainContent.addEventListener('click', (event) => {
        const postLink = event.target.closest('.post-title-link');
        const commentIcon = event.target.closest('.comment-icon');

        if (postLink || commentIcon) {
            const postId = postLink ? postLink.dataset.id : commentIcon.dataset.id;

            fetch(`/api/posts/${postId}/comments`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`Failed to fetch comments: ${response.status} ${response.statusText}`);
                    }
                    return response.json();
                })
                .then(comments => {
                    renderComments(comments);
                    dialog.classList.add('show');
                    dialogOverlay.classList.add('show');
                })
                .catch(error => {
                    console.error('Error fetching comments:', error);
                });
        }
    });

    dialogOverlay.addEventListener('click', () => {
        dialog.classList.remove('show');
        dialogOverlay.classList.remove('show');
    });
}
