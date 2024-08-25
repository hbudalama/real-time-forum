function initializeLikeDislikeButtons() {
    document.querySelectorAll('.like-button').forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');

            fetch(`/api/posts/${postId}/like`, {
                method: 'POST',
                credentials: 'include'
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        console.error('Server error:', text);
                        throw new Error(`Server error: ${response.status} ${response.statusText}`);
                    });
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    button.innerHTML = `<i class="fa fa-thumbs-up icon"></i>${data.likes}`;
                } else {
                    console.error('Error:', data.message);
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });

    document.querySelectorAll('.dislike-button').forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');

            fetch(`/api/posts/${postId}/dislike`, {
                method: 'POST',
                credentials: 'include'
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        console.error('Server error:', text);
                        throw new Error(`Server error: ${response.status} ${response.statusText}`);
                    });
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    button.innerHTML = `<i class="fa fa-thumbs-down icon"></i>${data.dislikes}`;
                } else {
                    console.error('Error:', data.message);
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });
}
