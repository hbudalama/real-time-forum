/**
 * Factory function to create a post element.
 * 
 * @param {Object} post - The post data.
 * @param {string} post.Username - The username of the post author.
 * @param {string} post.Title - The title of the post.
 * @param {string} post.Category - The category of the post.
 * @param {number} post.ID - The unique identifier for the post.
 * @param {number} post.Likes - The number of likes for the post.
 * @param {number} post.Dislikes - The number of dislikes for the post.
 * @param {number} post.Comments - The number of comments on the post.
 * @returns {HTMLDivElement} The HTMLDivElement representing the post.
 */
export function createPostFactory(post) {
    // Create the main post div
    const postDiv = document.createElement('div');
    postDiv.className = 'post';

    // Construct the inner HTML structure
    postDiv.innerHTML = `
        <div class="post-row">
            <div class="user-profile">
                <img src="/static/images/user.png" alt="${post.Username}'s profile picture">
                <div>
                    <p>${post.Username}</p>
                </div>
            </div>
        </div>
        <a href="javascript:void(0)" class="post-title-link" data-id="${post.ID}">
            <div>
                <h2>${post.Title}</h2>
                <h4>Category: ${post.Category}</h4>
            </div>
        </a>
        <div class="post-row">
            <div class="activity-icons">
                <div class="like-button" data-id="${post.ID}">
                    <i class="fa fa-thumbs-up icon"></i>${post.Likes}
                </div>
                <div class="dislike-button" data-id="${post.ID}">
                    <i class="fa fa-thumbs-down icon"></i>${post.Dislikes}
                </div>
                <div>
                    <a href="javascript:void(0)" class="comment-icon" data-id="${post.ID}">
                        <i class="fa fa-comment icon"></i>
                        <span id="post-${post.ID}-comments-count">${post.Comments}</span>
                    </a>
                </div>
            </div>
            <div class="post-profile-icon"></div>
        </div>
    `;

    return postDiv;
}