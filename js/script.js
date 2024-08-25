document.addEventListener('DOMContentLoaded', function () {
    fetch('/api/check_session')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Server error: ${response.status} ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.isAuthenticated) {
                loadForum();
            } else {
                loadLoginForm();
            }
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
            // Optionally display an error message to the user
        });
});

function loadLoginForm() {
    // Hide the navigation bar and left sidebar by adding a 'hidden' class
    document.querySelector('nav').classList.add('hidden');
    document.querySelector('.left-sidebar').classList.add('hidden');

    const formHtml = `
        <div id="login-form-container">
            <div id="form">
                <div id="button-container">
                    <input type="radio" id="loginBtn" name="btn" value="login" checked>
                    <label for="loginBtn">Login</label>
                    <input type="radio" id="registerBtn" name="btn" value="register">
                    <label for="registerBtn">Register</label>
                </div>
                <form id="loginField" class="input-group" action="/api/login" method="post" onsubmit="handleLogin(event)">
                    <div id="loginError" class="error-message"></div>
                    <input type="text" class="input-field" placeholder="Username or Email" name="username" required>
                    <input type="password" class="input-field" placeholder="Password" name="password" required>
                    <button type="submit" class="submit-btn">Log in</button>
                </form>
                <form id="registerField" style="display: none;" class="input-group" action="/api/signup" method="post" onsubmit="handleSignup(event)">
                    <div id="registerError" class="error-message"></div>
                    <input type="text" class="input-field" placeholder="Username" name="username" required>
                    <input type="text" class="input-field" placeholder="First name" name="firstName" required>
                    <input type="text" class="input-field" placeholder="Last name" name="lastName" required>
                    <label for="gender"> select your gender:</label>
                    <select class="gender" id="gender" name="gender">
                        <option value="1">female</option>
                        <option value="0">male</option>
                    </select>
                    <legend>select your age:</legend>
                    <select class="selectAge" id="selectAge" name="age">
                        <!-- Ages will be populated by server-side rendering -->
                    </select>
                    <input type="email" class="input-field" placeholder="Email" name="email" required>
                    <input type="password" class="input-field" placeholder="Password" name="password" required>
                    <input type="password" class="input-field" placeholder="Confirm Password" name="confirmPassword" required>
                    <button type="submit" class="submit-btn">Register</button>
                </form>
            </div>
        </div>
    `;
    document.getElementById('main-content').innerHTML = formHtml;
    loadAges();

    // Now that the elements exist, add event listeners
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const loginField = document.getElementById('loginField');
    const registerField = document.getElementById('registerField');
    const loginError = document.getElementById('loginError');
    const registerError = document.getElementById('registerError');

    loginBtn.addEventListener('click', () => {
        loginField.style.display = 'flex';
        registerField.style.display = 'none';
        loginError.style.display = 'none';
        registerError.style.display = 'none';
    });

    registerBtn.addEventListener('click', () => {
        loginField.style.display = 'none';
        registerField.style.display = 'flex';
        loginError.style.display = 'none';
        registerError.style.display = 'none';
    });
}

function loadForum() {
    // Show the navigation bar and left sidebar by removing the 'hidden' class
    document.querySelector('nav').classList.remove('hidden');
    document.querySelector('.left-sidebar').classList.remove('hidden');

    fetch('/api/posts')
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status} ${response.statusText}`);
        }
        return response.json();
    })
    .then(posts => {
        console.log(posts); // Inspect the data structure here
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
                                <i class="fa fa-comment icon"></i>${post.Comments}
                            </a>
                        </div>
                    </div>
                    <div class="post-profile-icon"></div>
                </div>
            </div>
        `).join('');

        document.getElementById('main-content').innerHTML = `<div class="index">${forumHtml}</div>`;
        // Reinitialize your forum JavaScript here
        initializePosts();
        initializeComments();
        initializeChat();
        initializeLikeDislikeButtons();
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });

}


function loadAges() {
    fetch('/api/get_ages')
        .then(response => response.json())
        .then(data => {
            const selectAge = document.getElementById('selectAge');
            data.ages.forEach(age => {
                const option = document.createElement('option');
                option.value = age;
                option.textContent = age;
                selectAge.appendChild(option);
            });
        });
}

function handleLogin(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    fetch('/api/login', {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
    .then(response => response.json())
    .then(data => {
        if (data.isAuthenticated) {
            loadForum();
        } else {
            document.getElementById('loginError').innerText = "Invalid login credentials.";
        }
    })
    .catch(error => console.error('Error:', error));
}

function logoutHandler() {
    fetch('/api/logout', {
        method: 'DELETE',
        credentials: 'include',

        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status} ${response.statusText}`);
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            loadLoginForm();
        } else {
            console.error('Logout failed:', data.message);
        }
    })
    .catch(error => console.error('There has been a problem with your fetch operation:', error));
}

function validateForm() {
    var checkboxes = document.querySelectorAll('input[name="post-category"]:checked');
    if (checkboxes.length === 0) {
        alert("Please select at least one category.");
        return false;
    }
    return true;
}

function handleSignup(event) {
    console.log("work u mf")
    event.preventDefault();
    const formData = new FormData(event.target);

    fetch("/api/signup", {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
    .then(response => {
        return response.json(); // Parse the JSON from the response
    })
    .then(data => {
        if (data.success) {
            alert("please log in.")
        } else {
            document.getElementById('registerField').innerText = data.reason || "Signup failed.";
        }
    })
    .catch(error => {
        document.getElementById('registerField').innerText = error.message || "An error occurred during signup.";
        console.error('Error:', error);
    });
}