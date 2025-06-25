import { globals, eventListeners } from './globals.js';
import { renderUserList, addEventFilter, addEventToNewCategory } from './render_home.js';
import { displayPostDetail } from './display_post_detail.js';
import { formatTime } from './format_time.js'
import { insertChat } from './display_chatroom.js';
export { openWebSocket };

function openWebSocket() {
    globals.ws = new WebSocket(`ws://localhost:8080/ws?user=${globals.LoggedInUser}`);

    globals.ws.onopen = function () {
        console.log("WebSocket connection established for " + globals.LoggedInUser);
        window.addEventListener('beforeunload', function () {
            globals.ws.close();
        });
    };

    globals.ws.onmessage = async function (event) {

        const message = JSON.parse(event.data);

        switch (message.type) {
            case "newUser":
                document.getElementById("user_list").innerHTML = ``;
                renderUserList();
                break;

            case "statusUpdate":
                updateUserStatus(message);
                break;

            case "categoryUpdate":

                // add category to the filter
                const filterTemplate = (categoryName) => `
                    <option value="${categoryName}">${categoryName}</option>`;
                const filterHTML = filterTemplate(message.content);
                document.getElementById("filter_list").insertAdjacentHTML('beforeend', filterHTML);

                // add category to post form
                const postFormTemplate = (category) => `
                    <label for="${category}">
                        <input type="checkbox" id="${category}" name="category" value="${category}">${category}
                    </label>`;
                const formHTML = postFormTemplate(message.content);
                document.getElementById("post_form").querySelector('div').insertAdjacentHTML('beforeend', formHTML);

                break;

            case "postUpdate":

                if (globals.currentFilter !== "all" && (!message.content.categories || !message.content.categories.includes(globals.currentFilter))) {
                    return;
                }

                const noPost = document.getElementById("no_post");
                if (noPost) {
                    noPost.remove();
                }

                const postTemplate = (id, title, author, date, numberOfComments) => `
                    <div id="post_${id}" class="post">
                        <h1 class="title">${title}</h1>
                        <h2 class="author">by ${author}</h2>
                        <h2 class="date">${date}</h2>
                        <div class="categories"></div>
                        <div class="numberOfComments">
                            <span id="number">${numberOfComments}</span>
                            <img src="../static/assets/chat_bubble.png" alt="chat bubble" />
                        </div>
                    </div>`;

                const postHTML = postTemplate(message.content.id, message.content.title, message.content.author, formatTime(message.content.time), message.content.numberOfComments);
                document.getElementById("posts_list").insertAdjacentHTML("afterbegin", postHTML);

                // insert categories
                console.log(message.content);
                const categoryTemplate = (categoryName) => `
                    <span class="category" >${categoryName}</span>`;
                if (message.content.categories) {
                    message.content.categories.forEach((category) => {
                        document.getElementById(`post_${message.content.id}`).querySelector('.categories').insertAdjacentHTML('beforeend', categoryTemplate(category));
                    })
                }

                // add event to open detail page
                document.getElementById(`post_${message.content.id}`).addEventListener('click', () => {
                    displayPostDetail(message.content);
                });

                break;

            case "privateMessage":

                // update sender's position in the list
                const senderUserDiv = document.getElementById(`user_${message.sender}`);
                const receiverUserDiv = document.getElementById(`user_${message.receiver}`);
                const userDiv = senderUserDiv ? senderUserDiv : receiverUserDiv;
                const userList = document.getElementById("user_list");
                userList.removeChild(userDiv);
                userList.insertAdjacentElement('afterbegin', userDiv);

                console.log(globals.currentPage, message.sender, message.receiver);
                console.log(globals.currentPage == message.sender);
                console.log(globals.LoggedInUser == message.receiver);


                // check the chat room is opened
                if (globals.currentPage == "chatRoom" && globals.chatRoomUser == message.sender || globals.chatRoomUser == message.receiver) {
                    console.log("opened");

                    // update chat room
                    const chatHistoryDOM = document.getElementById("chat_history")
                    if (chatHistoryDOM.innerHTML == "There is nothing here yet!") {
                        chatHistoryDOM.innerHTML = "";
                    }
                    insertChat(message.content, "new");
                    chatHistoryDOM.scrollTop = chatHistoryDOM.scrollHeight;
                    return;
                }

                // flag unread
                console.log(message.sender, globals.LoggedInUser);
                document.getElementById(`user_${message.sender}`).classList.add("unread");

                break;

            case "commentUpdate":
                // update comment number on post list

                console.log(globals.currentPage, message, globals.postID);
                if (globals.currentPage == "home") {
                    const postNumber = document.getElementById(`post_${message.content.postId}`).getElementsByClassName("number")[0];
                    postNumber.textContent = parseInt(postNumber.textContent) + 1;
                }

                // update comment on post detail
                console.log(globals.currentPage == "post", globals.postID == message.content.postId);
                if (globals.currentPage == "post" && globals.postID == message.content.postId) {
                    const number = document.getElementById("number");
                    number.textContent = parseInt(number.textContent) + 1;

                    const commentsDOM = document.getElementById("comments").querySelector("p");
                    if (commentsDOM.innerHTML == "There is nothing here yet!") {
                        commentsDOM.innerHTML = "";
                    }
                    const commentTemplate = (author, time, content) => `
                        <div class="comment">
                            <span class="commentName">${author}</span>
                            <span class="commentName">${time}</span>
                            <p>${content}</p>
                        </div>`;
                    document.getElementById(`comments`).insertAdjacentHTML('afterbegin', commentTemplate(message.content.author, formatTime(message.content.time), message.content.content));
                    document.getElementById(`comments`).scrollTo(0, 0)
                }

                break;
        }
    };

    globals.ws.onclose = function () {
        console.log("WebSocket connection closed");
    };

    globals.ws.onerror = function (error) {
        console.error("WebSocket error: ", error);
    };
}

function updateUserStatus(message) {
    const userElement = document.getElementById(`user_${message.sender}`);
    if (userElement) {
        const statusElement = userElement.querySelector("p");
        if (statusElement) {
            statusElement.textContent = message.content;
        }

        const statusIcon = userElement.querySelector("span");
        if (statusIcon) {
            statusIcon.className = message.content;
        }
    }

    if (globals.currentPage == "chatRoom" && globals.chatRoomUser == message.sender) {
        const chatroomStatus = document.getElementById(`chat_room_main`);
        chatroomStatus.querySelector("span").className = message.content;
        chatroomStatus.querySelector("p").textContent = message.content;
    }
}
