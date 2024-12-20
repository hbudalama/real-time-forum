import {initializeWebSocket} from './webSocket.js'
import { initializePosts } from './posts.js';
import { initializeLikeDislikeButtons } from './likesDislikes.js';
import { activeChatRecipient, chatClosed } from './chat.js';

export let loggedInUsername = null;

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
                initializeWebSocket();
                loadForum();
            } else {
                loadLoginForm();
            }
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });

    document.querySelector('.logo').addEventListener('click', function (event) {
        if (activeChatRecipient) {
            chatClosed()
        } 
        event.preventDefault();
        loadForum();
        console.log("didn't refresh")
    })
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
                <form id="loginField" class="input-group" action="/api/login" method="post" onsubmit="window.handleLogin(event)">
                    <div id="loginError" class="error-message"></div>
                    <input type="text" class="input-field" placeholder="Username or Email" name="username" required>
                    <input type="password" class="input-field" placeholder="Password" name="password" required>
                    <button type="submit" class="submit-btn">Log in</button>
                </form>
                <form id="registerField" style="display: none;" class="input-group" action="/api/signup" method="post" onsubmit="window.handleSignup(event)">
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
    const submitBtn = document.querySelector('.submit-btn');

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
    // document.getElementById('login-form-container').classList.add('hidden');
    // document.getElementById('form').classList.add('hidden');

    fetch('/api/get_user_info')
    .then(response => response.json())
    .then(data => {
        if (data.username) {
            loggedInUsername = data.username;
            const greetingDiv = document.getElementById('greeting');
            greetingDiv.textContent = `Hello, ${data.username}!`;
        }
    })
    .catch(error => {
        console.error('Error fetching user info:', error);
    });

    fetch('/api/posts')
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status} ${response.statusText}`);
        }
        return response.json();
    })
    .then(posts => {
        console.log(posts); // Inspect the data structure here
        
        // Check if there are any posts
        const mainContent = document.getElementById('main-content');
        if (posts && posts.length > 0) {
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

            mainContent.innerHTML = `<div class="index">${forumHtml}</div>`;
            // Hide registration form if there are posts
            // document.getElementById('form').classList.add('hidden');
        } else {
            mainContent.innerHTML = '<p style="color:white; font-size:1.5em; text-align:center;">No posts made yet. Be the First to create a post!</p>';
        }

        initializePosts();
        initializeComments();
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

window.handleLogin =  function handleLogin(event) {
    console.log(event)
    event.preventDefault()

    const formData = new FormData(event.target);
    fetch('/api/login', {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
    .then( async response => {
        if (!response.ok) {
            const  resp = await response.json()
            throw new Error(resp.reason)
        }
        return await response.json();
    })
    .then(data => {
        console.log("we got doata", data)
        if (data.isAuthenticated) {
            initializeWebSocket();
            loadForum();
        } else {
            const errorMessage = "Invalid login credentials.";
            console.error('Error:', errorMessage);
            // alert(`Error: ${errorMessage}`);
            document.getElementById('loginError').innerText = errorMessage;
        }
    })
    .catch(error => {
        console.error('Error:', error.message);
        alert(`Error: ${error.message}`);
    });
}


window.handleSignup = function handleSignup(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const password = formData.get('password');
    const confirmPassword = formData.get('confirmPassword');

    // Check if passwords match before sending the request
    if (password !== confirmPassword) {
        const errorMessage = "Passwords do not match!";
        console.error(errorMessage);
        alert(errorMessage);
        return; // Prevent form submission if passwords don't match
    }

    fetch("/api/signup", {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(errorData => {
                const errorMessage = errorData.reason || "Signup failed.";
                console.error('Error:', errorMessage);
                alert(`Error: ${errorMessage}`);
                throw new Error(errorMessage);
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            alert("Signup successful! Please log in.");
            // Optionally, redirect to the login page or clear the form here
        }
    })
    .catch(error => {
        console.error('Error:', error.message);
        // alert(`Error: ${error.message}`);
    });
}

