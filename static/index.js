// document.addEventListener('DOMContentLoaded', function () {
//     const filterSection = document.querySelector('.filter-section');
//     filterSection.addEventListener('change', applyFilters);

//     function applyFilters() {
//         const selectedCategories = [];
//         const selectedCriteria = {};

//         filterSection.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
//             if (checkbox.checked) {
//                 if (checkbox.name === 'filter-category') {
//                     selectedCategories.push(checkbox.value);
//                 } else {
//                     selectedCriteria[checkbox.value] = true;
//                 }
//             }
//         });

//         fetch('/filter-posts', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({
//                 categories: selectedCategories,
//                 criteria: selectedCriteria
//             }),
//         })
//             .then(response => response.json())
//             .then(data => {
//                 const mainContent = document.querySelector('.main-content');
//                 mainContent.innerHTML = ''; // Clear current posts

//                 data.posts.forEach(post => {
//                     const postElement = document.createElement('div');
//                     postElement.classList.add('post');
//                     postElement.innerHTML = `
//                     <div class="post-row">
//                         <div class="user-profile">
//                             <img src="/static/images/user.png">
//                             <div>
//                                 <p>${post.Username}</p>
//                             </div>
//                         </div>
//                     </div>
//                     <a href="/posts/${post.ID}" class="post-title-link">
//                         <div>
//                             <h2>${post.Title}</h2>
//                         </div>
//                     </a>
//                     <div class="post-row">
//                         <div class="activity-icons">
//                             <div>
//                                 <a href="/api/posts/${post.ID}/like">
//                                     <i class="fa fa-thumbs-up icon"></i>${post.Likes}
//                                 </a>
//                             </div>
//                             <div>
//                                 <a href="/api/posts/${post.ID}/dislike">
//                                     <i class="fa fa-thumbs-down icon"></i>${post.Dislikes}
//                                 </a>
//                             </div>
//                             <div>
//                                 <a href="/posts/${post.ID}">
//                                     <i class="fa fa-comment icon"></i>${post.Comments}
//                                 </a>
//                             </div>
//                         </div>
//                         <div class="post-profile-icon"></div>
//                     </div>
//                 `;
//                     mainContent.appendChild(postElement);
//                 });
//             })
//             .catch(error => console.log('Error fetching filtered posts:', error));
//     }
// });

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

   

    //Event listener for newest filter
    // document.getElementById('newest').addEventListener('change', function () {
    //     if (this.checked) {
    //         window.location.href = '/newest';
    //     } else {
    //         window.location.href = '/';
    //     }
    // });
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
