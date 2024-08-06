document.addEventListener('DOMContentLoaded', function () {
    var myPostsCheckbox = document.getElementById('myPosts');
    if (myPostsCheckbox) {
        myPostsCheckbox.addEventListener('change', function () {
            if (this.checked) {
                window.location.href = '/myPosts';
            } else {
                window.location.href = '/';
            }
        });
    }

    var myLikedPostsCheckbox = document.getElementById('Mylikedposts');
    if (myLikedPostsCheckbox) {
        myLikedPostsCheckbox.addEventListener('change', function () {
            if (this.checked) {
                window.location.href = '/Mylikedposts';
            } else {
                window.location.href = '/';
            }
        });
    }
});

function validateForm() {
    var checkboxes = document.querySelectorAll('input[name="post-category"]:checked');
    if (checkboxes.length === 0) {
        alert("Please select at least one category.");
        return false;
    }
    return true;
}

function logoutHandler(e) {
    fetch('/logout', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
        .then(response => {
            window.location.href = ('/login')
        })
}
