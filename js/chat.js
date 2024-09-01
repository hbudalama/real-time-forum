// function initializeChat() {
//     const chatUser = document.querySelector('.users-list');
//     const container = document.getElementById('main-content');
//     const nicknameElement = document.querySelector('#username p')
//     const userpic = document.querySelector('.user-icon img');

//     if (!nicknameElement) {
//         console.error("Nickname element is missing.");
//         return; 
//     }

//     const nickname = nicknameElement.textContent;

//     chatUser.addEventListener('click', () => {
//         const picsrc = userpic.src;
//         container.innerHTML = `<div id="chat-div">
//             <div id="user-info-chat">
//                 <img src="${picsrc}" id="user-pic">
//                 <p id="user-name-chat">${nickname}</p>
//             </div>
//         </div>`;
//     });
// }

function initializeChat(event) {
    console.log('open chat for', event.currentTarget.dataset.username)
}
