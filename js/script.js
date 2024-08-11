document.addEventListener('DOMContentLoaded', () => {
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

    loginField.addEventListener('submit', (e) => {
        e.preventDefault();
        const username = loginField.querySelector('input[name="username"]').value;
        const password = loginField.querySelector('input[name="password"]').value;
        loginError.style.display = 'none';

        if (username.includes(' ')) {
            loginError.textContent = "Username cannot contain spaces.";
            loginError.style.display = 'block';
            return;
        }
        if (password.length < 8 || password.includes(' ')) {
            loginError.textContent = "Password must be at least 8 characters long and cannot contain spaces.";
            loginError.style.display = 'block';
            return;
        }

        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        })
        .then(response => {
            if (response.ok) {
                window.location.href = '/';
            } else {
                return response.json().then(data => {
                    throw new Error(data.reason || "Unknown error");
                });
            }
        })
        .catch(error => {
            loginError.textContent = "An error occurred: " + error.message;
            loginError.style.display = 'block';
        });
    });

     document.getElementById("registerField").addEventListener("submit", function(event) {
            event.preventDefault();

            const formData = new FormData(this);
            const jsonObject = {};
            formData.forEach((value, key) => jsonObject[key] = value);

            fetch("/signup", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(jsonObject)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert("Registration successful!");
                } else {
                    alert("Registration failed: " + data.reason);
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
        });

    registerField.addEventListener('submit', (e) => {
        e.preventDefault();
        const username = registerField.querySelector('input[name="username"]').value;
        const firstName = registerField.querySelector('input[name="firstName"]').value;
        const lastName = registerField.querySelector('input[name="lastName"]').value;
        const gender = registerField.querySelector('select[name="gender"]').value;
        const age = registerField.querySelector('select[name="age"]').value;
        const email = registerField.querySelector('input[name="email"]').value;
        const password = registerField.querySelector('input[name="password"]').value;
        const confirmPassword = registerField.querySelector('input[name="confirmPassword"]').value;
        registerError.style.display = 'none';

        if (username.includes(' ')) {
            registerError.textContent = "Username cannot contain spaces.";
            registerError.style.display = 'block';
            return;
        }
        if (password.length < 8 || password.includes(' ')) {
            registerError.textContent = "Password must be at least 8 characters long and cannot contain spaces.";
            registerError.style.display = 'block';
            return;
        }
        if (password !== confirmPassword) {
            registerError.textContent = "Passwords do not match.";
            registerError.style.display = 'block';
            return;
        }

        const requestData = { username, firstName, lastName, gender, age, email, password, confirmPassword };
        console.log("Signup request data:", requestData);

        fetch('/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        })
        .then(response => {
            if (response.ok) {
                alert("Please log in");
                window.location.href = '/login';
            } else {
                return response.json().then(data => {
                    throw new Error(data.reason || "Unknown error");
                });
            }
        })
        .catch(error => {
            registerError.textContent = "An error occurred: " + error.message;
            registerError.style.display = 'block';
        });
    });
});

document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/check_session')
        .then(response => response.json())
        .then(data => {
            if (data.isAuthenticated) {
                loadForum();
            } else {
                loadLoginForm();
            }
        });
});

function loadLoginForm() {
    // Hide the navigation bar and left sidebar by adding a 'hidden' class
    document.querySelector('nav').classList.add('hidden');
    document.querySelector('.left-sidebar').classList.add('hidden');

    const formHtml = `
        <div id="form">
            <div id="button-container">
                <input type="radio" id="loginBtn" name="btn" value="login" checked>
                <label for="loginBtn">Login</label>
                <input type="radio" id="registerBtn" name="btn" value="register">
                <label for="registerBtn">Register</label>
            </div>
            <form id="loginField" class="input-group" action="/login" method="post">
                <div id="loginError" class="error-message"></div>
                <input type="text" class="input-field" placeholder="Username or Email" name="username" required>
                <input type="password" class="input-field" placeholder="Password" name="password" required>
                <button type="submit" class="submit-btn">Log in</button>
            </form>
            <form id="registerField" style="display: none;" class="input-group" action="/signup" method="post">
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

    const forumHtml = `
        <div class="index">
            <!-- Forum content here -->
        </div>
    `;
    document.getElementById('main-content').innerHTML = forumHtml;
    // Reinitialize your forum JavaScript here
    initializePosts();
    initializeComments();
    initializeChat();
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
