const loginPage = ```
    <div id="form">
        <div id="button-container">
            <input type="radio" id="loginBtn" name="btn" value="login" checked>
            <label for="loginBtn">Login</label>
            <input type="radio" id="registerBtn" name="btn" value="register">
            <label for="registerBtn">Register</label>
        </div>
        <form id="loginField" class="input-group" action="/login" method="post">
            <div id="loginError" class="error-message"></div>
            <input type="text" class="input-field" placeholder="Username" name="username" required>
            <input type="password" class="input-field" placeholder="Password" name="password" required>
            <button type="submit" class="submit-btn">Log in</button>
        </form>
        <form id="registerField" style="display: none;" class="input-group" action="/signup" method="post">
            <div id="registerError" class="error-message"></div>
            <input type="text" class="input-field" placeholder="Username" name="username" required>
            <input type="email" class="input-field" placeholder="Email" name="email" required>
            <input type="password" class="input-field" placeholder="Password" name="password" required>
            <input type="password" class="input-field" placeholder="Confirm Password" name="confirmPassword" required>
            <button type="submit" class="submit-btn">Register</button>
        </form>
    </div>
```


const payload = {
    usernmae: 'asda',
    password: 'asdasda'
}

fetch('/login', {
    body: payload,
    method: 'POST'
})
.then(async (res) => {
    if (!res.ok) {
        const body = await res.json()
        throw new Error('Invalid creditials: ' + body.reason)
    }
    hideLogin()
    showAddPost()
    showLoginItems()

})


function showAddPost() {
    document.getElementById('new-post-container').classList.add('hidden')
}


/**
 * 
 * @param {String} title 
 * @param {Function} func 
 * @returns {HTMLButtonElement}
 */
function createDropdownOption(title, func) {
    const item = document.createElement('button')
    item.classList.add('drop-item')
    item.innerText = title
    item.onclick = func
    
}

function better() {
    const dropdown = document.getElementById('user-dropdown-items')

    dropdown.appendChild(
        createDropdownOption('Logout', handleLogout),
        createDropdownOption('Profile', showProfile)
    )

}


function showLoginItems() {
    const dropdown = document.getElementById('user-dropdown-items')
    // dropdown.innerHTML = ''

    // const options = ```
    // <div class="dropdown-item" onclick="handleLogout()">Logout</div>
    // <div class="dropdown-item">Profile</div>
    // ```

    // dropdown.innerHTML = options


    // const loggedDiv = document.getElementById('user-logged-options');
    // loggedDiv.attributes.removeNamedItem('hidden')

    // const loggedDiv = document.getElementById('user-logged-options');


}