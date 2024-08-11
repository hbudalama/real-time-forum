function initializeChat() {
    const chatUser = document.querySelector('.users-list');
    const container = document.getElementById('main-content');
    const nickname = document.querySelector('#username p').textContent;
    const userpic = document.querySelector('.user-icon img');

    chatUser.addEventListener('click', () => {
        const picsrc = userpic.src;
        container.innerHTML = `<div id="chat-div">
            <div id="user-info-chat">
                <img src="${picsrc}" id="user-pic">
                <p id="user-name-chat">${nickname}</p>
            </div>
        </div>`;
    });
}
