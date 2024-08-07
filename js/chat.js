document.addEventListener('DOMContentLoaded', () => {
    const chatUser = document.querySelector('.users-list');
    const container = document.getElementById('main-content');
    const nickname = document.querySelector('#username p').textContent;
    const userpic = document.querySelector('.user-icon img');

    chatUser.addEventListener('click', () => {
        const picsrc = userpic.src;
        container.innerHTML = `<div id="chat-div">
                <p><img src="${picsrc}" id="user-pic">${nickname}</p>
             </div>`
    })
});

