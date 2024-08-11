function initializeComments() {
    const mainContent = document.getElementById('main-content');
    const dialog = document.getElementById('comment-dialog');
    const dialogOverlay = document.getElementById('dialog-overlay');
    const commentList = document.getElementById('comment-list');

    const sampleComments = [
        { Username: 'Alice', Content: 'Great post!', Likes: 10, Dislikes: 2, ID: 1, CreatedDate: new Date() },
        { Username: 'Bob', Content: 'Thanks for sharing!', Likes: 5, Dislikes: 1, ID: 2, CreatedDate: new Date() },
    ];

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
        if (event.target.closest('.post-title-link') || event.target.closest('.comment-icon')) {
            renderComments(sampleComments);
            dialog.classList.add('show');
            dialogOverlay.classList.add('show');
        }
    });

    dialogOverlay.addEventListener('click', () => {
        dialog.classList.remove('show');
        dialogOverlay.classList.remove('show');
    });
}
