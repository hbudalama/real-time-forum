function initializePosts() {
    const postBtn = document.getElementById('add-post-btn');
    const container = document.getElementById('main-content');

    postBtn.addEventListener('click', () => {
        // Check if the post creation form is already visible
        const existingForm = document.getElementById('create-post-form');
        if (existingForm) {
            console.log('Post creation form is already visible.');
            return; // Exit if form is already visible
        }

        // Fetch the logged-in username
        fetch('/api/get_user_info')
            .then(response => response.json())
            .then(data => {
                const username = data.username || 'username'; // Fallback to 'username' if no data
                const pc = document.getElementById('postss-container')
                pc.classList.add('hidden')
                // container.innerHTML = ` 
                //     <div class="post-container">
                //         <div class="user-profile-post">
                //             <img src="/static/images/user.png">
                //             <div>   
                //                 <p>${username}</p>
                //             </div>
                //         </div>
                //         <form id="create-post-form">
                //             <textarea name="title" rows="3" placeholder="What's on your mind?"></textarea>
                //             <textarea name="content" rows="3" placeholder="Content..."></textarea>
                //             <p>Choose a category:</p>
                //             <ul>
                //                 <li>
                //                     <input type="checkbox" id="post-education" name="post-category" value="education">
                //                     <label for="post-education">Education</label>
                //                 </li>
                //                 <li>
                //                     <input type="checkbox" id="post-entertainment" name="post-category" value="entertainment">
                //                     <label for="post-entertainment">Entertainment</label>
                //                 </li>
                //                 <li>
                //                     <input type="checkbox" id="post-sports" name="post-category" value="sports">
                //                     <label for="post-sports">Sports</label>
                //                 </li>
                //                 <li>
                //                     <input type="checkbox" id="post-news" name="post-category" value="news">
                //                     <label for="post-news">News</label>
                //                 </li>
                //             </ul>
                //             <div id="post-button-container">
                //                 <button class="post-button" type="submit">Post</button>
                //             </div>
                //         </form>                    
                //     </div> 
                // `;

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
                                console.log("i have created a post!")
                                // Hide the form and reload posts
                                container.innerHTML = '';
                                loadPosts(); // Load and display posts dynamically

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
    fetch('/api/posts')
        .then(response => response.json())
        .then(data => {
            const postsContainer = document.getElementById('posts-container');
            if (!postsContainer) {
                console.error('Posts container element not found.');
                return;
            }
            postsContainer.innerHTML = ''; // Clear existing posts
            data.posts.forEach(post => {
                postsContainer.innerHTML += `
                    <div class="post">
                        <div class="user-profile-post">
                            <img src="/static/images/user.png">
                            <div>   
                                <p>${post.username}</p>
                            </div>
                        </div>
                        <h3>${post.title}</h3>
                        <p>${post.content}</p>
                        <p>Category: ${post.category}</p>
                    </div>
                `;
            });
        })
        .catch(error => console.error('Error fetching posts:', error));
}

document.addEventListener('DOMContentLoaded', initializePosts);