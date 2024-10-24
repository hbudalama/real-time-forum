import { activeChatRecipient, chatClosed } from './chat.js';
import { initializeLikeDislikeButtons } from './likesDislikes.js';

export function initializePosts() {
    const postBtn = document.getElementById('add-post-btn');
    const container = document.getElementById('main-content');

    postBtn.addEventListener('click', () => {
        if (activeChatRecipient) {
            chatClosed()
        }
        const existingForm = document.getElementById('create-post-form');
        if (existingForm) {
            console.log('Post creation form is already visible.');
            return;
        }

        fetch('/api/get_user_info')
            .then(response => response.json())
            .then(data => {
                const username = data.username || 'username';
                container.innerHTML = ` 
                    <div class="post-container">
                        <div class="user-profile-post">
                            <img src="/static/images/user.png">
                            <div>   
                                <p>${username}</p>
                            </div>
                        </div>
                        <form id="create-post-form">
                            <textarea name="title" rows="3" placeholder="What's on your mind?"></textarea>
                            <textarea name="content" rows="3" placeholder="Content..."></textarea>
                            <p>Choose a category:</p>
                            <ul>
                                <li>
                                    <input type="radio" id="post-education" name="post-category" value="education">
                                    <label for="post-education">Education</label>
                                </li>
                                <li>
                                    <input type="radio" id="post-entertainment" name="post-category" value="entertainment">
                                    <label for="post-entertainment">Entertainment</label>
                                </li>
                                <li>
                                    <input type="radio" id="post-sports" name="post-category" value="sports">
                                    <label for="post-sports">Sports</label>
                                </li>
                                <li>
                                    <input type="radio" id="post-news" name="post-category" value="news">
                                    <label for="post-news">News</label>
                                </li>
                            </ul>
                            <div id="post-button-container">
                                <button class="post-button" type="submit">Post</button>
                            </div>
                        </form>                    
                    </div> 
                `;

                const form = document.getElementById('create-post-form');
                form.addEventListener('submit', function (event) {
                    event.preventDefault();
                    // Call validateForm before submitting
                      if (!validateForm()) {
                        return; // If the form is invalid, stop the submission
                    }

                    const formData = new FormData(form);

                    fetch('/api/add-post', {
                        method: 'POST',
                        body: formData,
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                alert('Post created successfully!');
                                loadPosts(); // Reload posts
                                
                                // Re-initialize like and dislike buttons after new posts are loaded
                                initializeLikeDislikeButtons();
                            } else {
                                alert('Failed to create post.');
                            }
                        })
                        .catch(error => console.error('Error:', error));
                });
            })
            .catch(error => console.error('Error fetching user info:', error));
    });
}


// Function to load and display posts dynamically
function loadPosts() {
    const container = document.getElementById('main-content');
    fetch('/api/posts')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Server error: ${response.status} ${response.statusText}`);
            }
            return response.json();
        })
        .then(posts => {
            console.log('Posts data:', posts);

            if (!Array.isArray(posts)) {
                console.error('Expected an array of posts but received:', posts);
                return;
            }

            if (posts.length === 0) {
                console.log('No posts made.');
                return;
            }

            const forumHtml = posts.map(post => `
                <div class="post">
                    <div class="post-row">
                        <div class="user-profile">
                            <img src="/static/images/user.png">
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
                </div>
            `).join('');

            document.getElementById('main-content').innerHTML = `<div class="index">${forumHtml}</div>`;

            initializePosts();
            initializeComments();
            initializeLikeDislikeButtons();
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}

function validateForm() {
    const title = document.querySelector('textarea[name="title"]').value.trim();
    const content = document.querySelector('textarea[name="content"]').value.trim();
    const category = document.querySelector('input[name="post-category"]:checked');

    if (!title || !content || !category) {
        alert("The post must have a title, content, and category.");
        return false;
    }

    if (title.length > 100) {
        alert('Title is too long.');
        return;
    }

    if (content.length > 500) {
        alert('Content limit is 500 characters.');
        return;
    }

    return true;
}
