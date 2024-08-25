function initializePosts() {
    const postBtn = document.getElementById('add-post-btn');
    const container = document.getElementById('main-content');
    const posts = document.querySelector('.post');

    postBtn.addEventListener('click', () => {
        if (posts) {
            posts.style.display = 'none';
        } else {
            console.error('No posts element found.');
        }

        // Fetch the logged-in username
        fetch('/api/get_user_info')
            .then(response => response.json())
            .then(data => {
                const username = data.username || 'username'; // Fallback to 'username' if no data
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
                                    <input type="checkbox" id="post-education" name="post-category" value="education">
                                    <label for="post-education">Education</label>
                                </li>
                                <li>
                                    <input type="checkbox" id="post-entertainment" name="post-category" value="entertainment">
                                    <label for="post-entertainment">Entertainment</label>
                                </li>
                                <li>
                                    <input type="checkbox" id="post-sports" name="post-category" value="sports">
                                    <label for="post-sports">Sports</label>
                                </li>
                                <li>
                                    <input type="checkbox" id="post-news" name="post-category" value="news">
                                    <label for="post-news">News</label>
                                </li>
                            </ul>
                            <div id="post-button-container">
                                <button class="post-button" type="submit">Post</button>
                            </div>
                        </form>                    
                    </div> 
                `;

                // Handle form submission via fetch
                const form = document.getElementById('create-post-form');
                form.addEventListener('submit', function (event) {
                    event.preventDefault();
                    const formData = new FormData(form);
                    fetch('/api/add-post', {
                        method: 'POST',
                        body: formData,
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                alert('Post created successfully!');
                                // Reload posts or update the post list dynamically
                                loadPosts();
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
    fetch('/api/post')
        .then(response => response.json())
        .then(data => {
            // Code to dynamically display posts
        })
        .catch(error => console.error('Error fetching posts:', error));
}

document.addEventListener('DOMContentLoaded', initializePosts);