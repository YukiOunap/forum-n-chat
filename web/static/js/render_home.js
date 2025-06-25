import { globals, eventListeners } from './globals.js';
import { displayLogin } from './login_signup.js';
import { displayPostDetail } from './display_post_detail.js';
import { displayChatRoom } from './display_chatroom.js'
import { debounce } from './debounce.js'
import { openWebSocket } from './websocket.js'
import { formatTime } from './format_time.js'
export { renderHome, renderPosts, fetchPostCategories, addEventFilter, addEventToNewCategory, renderUserList };

function renderHome() {
    // setup header
    document.getElementById("nickname").textContent = globals.LoggedInUser;
    if (!eventListeners.addEventLogout) {
        addEventLogout();
        eventListeners.addEventLogout = true;
    }

    // setup window
    document.getElementById("filter_list").innerHTML = '<option value="all">No filter (All posts)</option>';
    renderFilter();
    document.getElementById("post_form").querySelector('div').innerHTML = '';
    renderPostForm();

    // initialize (clear) the posts
    document.getElementById("posts_list").innerHTML = ``;
    globals.postOffset = 0;
    renderPosts("all");

    // setup user list
    document.getElementById("user_list").innerHTML = '';
    renderUserList();

    // add event listener
    addEventPostScroll();

    // display home page
    globals.currentPage = "home";
    document.getElementById("login").style.display = 'none';
    document.getElementById("home").style.display = 'block';
    document.getElementById("signup").style.display = 'none';
    document.getElementById("create_post_form_bg").style.display = 'none';
    // inside the window
    document.getElementById('posts').style.display = 'block';
    document.getElementById('chat_room').style.display = 'none';
    document.getElementById('post_detail').style.display = 'none';

    openWebSocket();
}

function addEventLogout() {
    document.getElementById('logout').addEventListener('click', async function () {

        try {
            if (globals.ws) globals.ws.close();

            const response = await fetch('/logout', {
                method: 'POST',
            });
            if (response.ok) {
                globals.LoggedInUser = "";
                globals.currentPage = "login"
                displayLogin();
            } else {
                console.error("Logout failed.");
            }
        } catch (error) {
            console.error('Error logging out:', error);
        }
    });
}

async function renderFilter() {

    // insert all categories to filter
    const categories = await fetchCategories();

    if (categories == null) {
        return;
    }

    categories.forEach((category) => {
        const categoryTemplate = (categoryName) => `
            <option value="${categoryName}">${categoryName}</option>`;

        const categoryHTML = categoryTemplate(category.name);
        document.getElementById("filter_list").insertAdjacentHTML('beforeend', categoryHTML);
    });

    addEventFilter();
}

async function fetchCategories() {
    try {
        const response = await fetch(`/fetch-categories`);
        const data = await response.json();
        return data.categories;
    } catch (error) {
        console.error('Error fetching categories:', error);
        return [];
    }
}

function addEventFilter() {
    document.getElementById("filter_form").addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(event.target)
        globals.postOffset = 0;
        globals.currentFilter = formData.get('filter')

        document.getElementById("posts_list").innerHTML = ``;
        renderPosts(globals.currentFilter);
    });
}

async function renderPosts(filter) {
    const posts = await fetchPosts(filter);
    if (globals.postOffset == 0 && !posts) { // offset 0 = not a scroll event
        document.getElementById("posts_list").innerHTML = `
            <p id="no_post"> There is nothing here yet!</p>`;
        return;
    }
    globals.postOffset += 10;

    const postCategories = await fetchPostCategories();
    if (!posts) {
        return;
    }
    document.getElementById("posts_list")
    posts.forEach((post) => {
        renderPost(post, postCategories);
    })
}

function addEventPostScroll() {
    const postsWindow = document.getElementById("window");
    postsWindow.addEventListener('scroll', debounce(function () {

        if (globals.currentPage == "home" && postsWindow.scrollTop + postsWindow.clientHeight >= postsWindow.scrollHeight) {
            renderPosts(globals.currentFilter);
        }
    }, 200));
}

async function fetchPosts(filter) {
    try {
        const response = await fetch(`/fetch-posts?category=${filter}&offset=${globals.postOffset}`);
        const data = await response.json();
        return data.posts;
    } catch (error) {
        console.error('Error fetching posts:', error);
        return [];
    }
}

async function renderPost(post, postCategories) {

    const postTemplate = (id, title, author, date, numberOfComments) => `
        <div id="post_${id}" class="post">
            <h1 class="title">${title}</h1>
            <h2 class="author">by ${author}</h2>
            <h2 class="date">${date}</h2>
            <div class="categories"></div>
            <div class="numberOfComments">
                <span class="number">${numberOfComments}</span>
                <img src="../static/assets/chat_bubble.png" alt="chat bubble" />
            </div>
        </div>`;

    const postHTML = postTemplate(post.id, post.title, post.author, formatTime(post.time), post.numberOfComments);
    document.getElementById("posts_list").insertAdjacentHTML('beforeend', postHTML);

    // insert categories
    if (postCategories) {
        const categoryTemplate = (categoryName) => `
            <span class="category" >${categoryName}</span>`;
        postCategories.forEach((category) => {
            if (category.postId == post.id) {
                document.getElementById(`post_${post.id}`).querySelector('.categories').insertAdjacentHTML('beforeend', categoryTemplate(category.categoryName));
            }
        })
    }

    // add event to open detail page
    document.getElementById(`post_${post.id}`).addEventListener('click', () => {
        displayPostDetail(post);
    });
}

async function fetchPostCategories() {
    try {
        const response = await fetch(`/fetch-post-categories`);
        const data = await response.json();
        return data.postCategories;
    } catch (error) {
        console.error('Error fetching post categories:', error);
        return [];
    }
}

async function renderPostForm() {

    // add registered categories to the form
    const categories = await fetchCategories();

    if (categories != null) {
        const categoryTemplate = (category) => `
        <label for="${category}">
            <input type="checkbox" id="${category}" name="category" value="${category}">${category}
        </label>`;
        categories.forEach((category) => {
            const categoryHTML = categoryTemplate(category.name);
            document.getElementById("post_form").querySelector('div').insertAdjacentHTML('beforeend', categoryHTML);
        })
    }

    // add events
    if (!eventListeners.addEventToNewCategory) {
        addEventToNewCategory();
        eventListeners.addEventToNewCategory = true;
    }
    if (!eventListeners.addEventToCreatePost) {
        addEventToCreatePost();
        eventListeners.addEventToCreatePost = true;
    }

    // add event to open/close the form
    document.getElementById('create_post').addEventListener('click', function () {
        document.getElementById("create_post_form_bg").style.display = 'block';
        document.getElementById("create_post_form").style.display = 'block';
    })

    document.getElementById('create_post_form_bg').addEventListener('click', function (event) {
        if (document.getElementById("create_post_form").contains(event.target)) {
            return;
        }
        document.getElementById("create_post_form_bg").style.display = 'none';
        document.getElementById("create_post_form").style.display = 'none';
    })
    document.getElementById('create_post_form').querySelector('img').addEventListener('click', function () {
        document.getElementById("create_post_form_bg").style.display = 'none';
        document.getElementById("create_post_form").style.display = 'none';
    })
};

function addEventToNewCategory() {

    document.getElementById("add_category").addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(event.target);
        try {
            const response = await fetch('/add-category', {
                method: 'POST',
                body: formData
            });
            const data = await response.json();
            if (data.status == "duplicates") {
                alert("This category is already registered");
                return;
            }
            document.querySelector('#add_category input[name="newCategory"]').value = '';
        } catch (error) {
            console.error('Error creating post:', error);
        }
    });
}


function addEventToCreatePost() {
    document.getElementById("post_form").addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(event.target);
        formData.append("author", globals.LoggedInUser);
        try {
            await fetch('/create-post', {
                method: 'POST',
                body: formData
            });
            document.querySelectorAll('#post_form input[type="checkbox"]').forEach((checkbox) => {
                checkbox.checked = false;
            });
            document.getElementById("title_input").value = '';
            document.getElementById("content_input").value = '';

            document.getElementById("create_post_form_bg").style.display = 'none';
            document.getElementById("create_post_form").style.display = 'none';
        } catch (error) {
            console.error('Error creating post:', error);
        }
    })
}

async function renderUserList() {
    const userData = await fetchUserList();
    const latestMessages = await fetchLastMessages(globals.LoggedInUser);

    // Initialize an empty array to store sorted users
    let sortedUsers = [];
    const latestMessageMap = new Map();

    // If there are latestMessages, process them
    if (latestMessages != null) {
        // Create a map of latest message time by use
        latestMessages.forEach(message => {
            const user1 = message.user1;
            const user2 = message.user2;

            if (user1 !== globals.LoggedInUser) {
                latestMessageMap.set(user1, message);
            }
            if (user2 !== globals.LoggedInUser) {
                latestMessageMap.set(user2, message);
            }
        });
    }
    console.log(latestMessageMap);


    // Separate users with recent messages and those without
    const messagedUsers = [];
    const noMessageUsers = [];

    userData.users.forEach(user => {
        if (user !== globals.LoggedInUser) {
            if (latestMessageMap.has(user)) {
                messagedUsers.push(user);
            } else {
                noMessageUsers.push(user);
            }
        }
    });
    console.log(messagedUsers, noMessageUsers);


    // Sort users with messages by latest message time (descending)
    messagedUsers.sort((a, b) => {
        const timeA = latestMessageMap.get(a)?.lastMessageTime || 0; // Default to 0 if no message
        const timeB = latestMessageMap.get(b)?.lastMessageTime || 0; // Default to 0 if no message
        return new Date(timeB) - new Date(timeA); // Sort descending by time
    });

    // Sort users without messages alphabetically
    noMessageUsers.sort();

    console.log(messagedUsers, noMessageUsers);


    // Combine both lists
    sortedUsers = [...messagedUsers, ...noMessageUsers];

    // Insert users into the DOM
    for (const user of sortedUsers) {
        const isRead = await fetchIsRead(globals.LoggedInUser, user);
        console.log(isRead);

        if (userData.onlineUsers.includes(user)) {
            InsertUser(user, "online", isRead.isRead);
        } else {
            InsertUser(user, "offline", isRead.isRead);
        }
    }
}

async function fetchUserList() {
    try {
        const response = await fetch(`/fetch-users`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user data:', error);
        return [];
    }
}

async function fetchLastMessages(user) {
    try {
        const response = await fetch(`/fetch-latest-messages?user=${user}`);
        const data = await response.json();
        return data.latestMessages;
    } catch (error) {
        console.error('Error fetching user data:', error);
        return [];
    }
}

async function fetchIsRead(me, other) {
    try {
        const response = await fetch(`/fetch-is-read?me=${me}&other=${other}`);
        const data = await response.json();
        console.log(data);

        return data;
    } catch (error) {
        console.error('Error fetching user data:', error);
        return [];
    }
}

function InsertUser(user, status, isRead) {

    console.log(user, status, isRead);
    const userTemplate = (user, status) => `
    <div id="user_${user}" class="user">
        <h3 class="account-name">${user}</h3>
        <span class="${status}"></span>
        <p>${status}</p>
    </div>`;

    const onlineUserHTML = userTemplate(user, status);
    document.getElementById("user_list").insertAdjacentHTML('beforeend', onlineUserHTML);

    // add unread class for the users with unread message
    const insertedUser = document.getElementById(`user_${user}`)
    if (isRead !== null && isRead !== undefined && isRead === false) {
        insertedUser.classList.add('unread');
    }

    insertedUser.addEventListener('click', () => {
        displayChatRoom(user, status);

        // color only the selected user
        for (let user of document.getElementsByClassName('user')) {
            user.classList.remove('selected');
        };
        insertedUser.classList.add('selected');
    });
}