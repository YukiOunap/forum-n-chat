import { globals, eventListeners } from './globals.js';
import { debounce } from './debounce.js'
import { formatTime } from './format_time.js'
export { displayChatRoom, insertChat };

let offset = 0;

function displayChatRoom(user, status) {
    const chatRoomTemplate = (user, status) => `
        <div class="user">
            <h2 class="account-name">${user}</h2>
            <span class="${status}"></span>
            <p>${status}</p>
        </div>

        <div id="chat_history"></div>

        <form id="chat_form_${user}">
            <textarea type="text" id="message_input_${user}" placeholder="Type a message" required></textarea>
            <button class="button" type="submit">Send</button>
        </form>`;

    const chatRoomHTML = chatRoomTemplate(user, status);
    document.getElementById("chat_room_main").innerHTML = chatRoomHTML;

    // display chat history
    offset = 0;
    displayChatHistory(user);

    // add back home event
    document.getElementById(`chat_room`).getElementsByClassName(`back_home_link`)[0].addEventListener('click', () => {
        document.getElementById("posts").style.display = 'block';
        document.getElementById("chat_room").style.display = 'none';
    });

    // add send message event
    document.getElementById(`chat_form_${user}`).addEventListener('submit', function (event) {
        event.preventDefault();

        globals.ws.send(JSON.stringify({
            type: 'privateMessage',
            sender: globals.LoggedInUser,
            receiver: user,
            content: document.getElementById(`message_input_${user}`).value
        }));
        document.getElementById(`message_input_${user}`).value = "";
    });

    globals.currentPage = "chatRoom";
    globals.chatRoomUser = user;
    document.getElementById("posts").style.display = 'none';
    document.getElementById("chat_room").style.display = 'block';
    document.getElementById("post_detail").style.display = 'none';
}

async function displayChatHistory(user) {
    console.log(`/render-chat-history?sender=${globals.LoggedInUser}&receiver=${user}&offset=${offset}`)

    try {
        const response = await fetch(`/render-chat-history?sender=${globals.LoggedInUser}&receiver=${user}&offset=${offset}`);
        const messages = await response.json();

        const chatHistoryDOM = document.getElementById("chat_history");

        if (!messages && offset == 0) {
            chatHistoryDOM.innerHTML = "There is nothing here yet!";
            return;
        }

        const previousScrollHeight = chatHistoryDOM.scrollHeight;
        messages.forEach(message => {
            insertChat(message, "history");
        })
        if (offset == 0) {
            chatHistoryDOM.scrollTop = chatHistoryDOM.scrollHeight;
        } else {
            const updatedScrollHeight = chatHistoryDOM.scrollHeight;
            chatHistoryDOM.scrollTop = updatedScrollHeight - previousScrollHeight;
        }

        offset += 10;

        // add scroll event
        document.getElementById("chat_history").addEventListener('scroll', debounce(function () {
            if (chatHistoryDOM.scrollTop === 0) {
                displayChatHistory(user);
            }
        }, 200));

        const userDiv = document.getElementById(`user_${user}`);
        if (userDiv) {
            userDiv.classList.remove('unread');
        }
    } catch (error) {
        console.error('Error fetching messages:', error)
    }
}

function insertChat(message, chatType) {

    const chatHistoryDOM = document.getElementById("chat_history");

    let type = "me";
    if (message.receiver == globals.LoggedInUser) {
        type = "other";
    }
    const chatHistoryTemplate = (type, content, time) => `
                <div class="${type}">
                    <p>${content}</p>
                    <p class="time">${time}</p>
                </div>`;

    const chatHistoryHTML = chatHistoryTemplate(type, message.content, formatTime(message.time));
    if (chatType == "new") {
        chatHistoryDOM.insertAdjacentHTML('beforeend', chatHistoryHTML);
        return;
    }
    chatHistoryDOM.insertAdjacentHTML('afterbegin', chatHistoryHTML);
}