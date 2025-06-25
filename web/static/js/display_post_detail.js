import { globals, eventListeners } from './globals.js';
import { debounce } from './debounce.js'
import { fetchPostCategories, renderPosts } from './render_home.js';
import { formatTime } from './format_time.js'
export { displayPostDetail };

let offset = 0;

async function displayPostDetail(post) {
    eventListeners.addEventToCommentForm = false;
    document.getElementById("post_detail_main").innerHTML = '';

    const fetchedPost = await fetchPost(post.id);

    const postTemplate = (id, title, author, date, numberOfComments, content) => `
        <div id="post">
            <h1 class="title">${title}</h1>
            <h2 class="author">by ${author}</h2>
            <h2 class="date">${date}</h2>
            <div class="categories"></div>
            <div class="numberOfComments">
                <span id="number">${numberOfComments}</span>
                <img src="../static/assets/chat_bubble.png" alt="chat bubble" />
            </div>
            <p>${content}</p>
        </div>

        <h3>Comment</h2>
        <div id="comments">
            <p class="noComment">There is nothing here yet!</p>
        </div>
        
        <form id="comment_form_${id}">
            <textarea type="text" id="comment_input" name="content" placeholder="Type a comment" required></textarea>
            <button class="button" type="submit">Send</button>
        </form>`;

    const postHTML = postTemplate(fetchedPost.id, fetchedPost.title, fetchedPost.author, formatTime(fetchedPost.time), fetchedPost.numberOfComments, fetchedPost.content);
    document.getElementById("post_detail_main").innerHTML = postHTML;

    // insert category
    const postCategories = await fetchPostCategories();
    if (postCategories) {
        const categoryTemplate = (categoryName) => `
        <span class="category" >${categoryName}</span>`;
        postCategories.forEach((category) => {
            if (category.postId == post.id) {
                document.getElementById(`post_detail`).querySelector('.categories').insertAdjacentHTML('beforeend', categoryTemplate(category.categoryName));
            }
        })
    }

    // insert comments
    offset = 0;
    renderComments(post);

    // add event listener to back home
    if (!eventListeners.postBackHome) {
        document.getElementById(`post_detail`).getElementsByClassName(`back_home_link`)[0].addEventListener('click', function (event) {
            document.getElementById("posts_list").innerHTML = ``;
            globals.postOffset = 0;
            renderPosts(globals.currentFilter);

            document.getElementById("posts").style.display = 'block';
            document.getElementById("post_detail").style.display = 'none';
        });
        eventListeners.postBackHome = true;
    }

    // add event to comment form
    if (!eventListeners.addEventToCommentForm) {
        addEventToCommentForm(post);
        eventListeners.addEventToCommentForm = true;
    }

    // add event for load comments
    const comments = document.getElementById(`comments`);
    comments.addEventListener('scroll', debounce(function () {
        if (globals.currentPage == "post" && comments.scrollTop + comments.clientHeight >= comments.scrollHeight) {
            renderComments(post);
        }
    }, 200));

    globals.currentPage = "post";
    globals.postID = post.id
    document.getElementById("posts").style.display = 'none';
    document.getElementById("post_detail").style.display = 'block';
}

async function fetchPost(postId) {
    try {
        const response = await fetch(`/fetch-post?id=${postId}`);
        const data = await response.json();
        return data.post;
    } catch (error) {
        console.error('Error fetching posts:', error);
        return;
    }
}

async function renderComments(post) {
    const comments = await fetchComments(post.id);

    const emptyComment = document.getElementById(`comments`).getElementsByClassName("noComment")[0];
    if (offset != 0 || comments != null) {
        emptyComment.innerHTML = "";
    }
    const commentTemplate = (author, time, content) => `
        <div class="comment">
            <span class="commentName">${author}</span>
            <span class="commentName">${time}</span>
            <p>${content}</p>
        </div>`;
    comments.forEach((comment) => {
        document.getElementById(`comments`).insertAdjacentHTML('beforeend', commentTemplate(comment.author, formatTime(comment.time), comment.content));
    })

    offset += 10;
}

async function fetchComments(id) {
    try {
        const response = await fetch(`/fetch-comments?post=${id}&offset=${offset}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching comments:', error);
        return [];
    }
}

function addEventToCommentForm(post) {
    document.getElementById(`comment_form_${post.id}`).addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(event.target);
        formData.append("author", globals.LoggedInUser);
        formData.append("postId", post.id);
        try {
            await fetch('/create-comment', {
                method: 'POST',
                body: formData
            });
            document.getElementById("comment_input").value = '';
        } catch (error) {
            console.error('Error creating post:', error);
        }
    })
}