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
        
        container.innerHTML = ` 
            <div class="post-container">
                <div class="user-profile-post">
                    <img src="/static/images/user.png">
                    <div>
                        <p>username</p>
                    </div>
                </div>
                {{if .ErrorMessage}}
                <div class="error-message">{{.ErrorMessage}}</div>
                {{end}}
                <form action="/add-post" method="POST" onsubmit="return validateForm()">
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
    });
}

document.addEventListener('DOMContentLoaded', initializePosts);
