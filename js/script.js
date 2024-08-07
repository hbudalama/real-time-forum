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
