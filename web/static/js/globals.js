// global variables
const globals = {
    currentPage: "login",
    chatRoomUser: "",
    postID: "",
    LoggedInUser: "",
    ws: "",
    currentFilter: "all",
    postOffset: 0,
};

const eventListeners = {
    addEventSignupLink: false,
    addEventToLoginForm: false,
    addEventLoginLink: false,
    addEventToSignUpForm: false,

    addEventLogout: false,
    addEventToNewCategory: false,
    addEventToCreatePost: false,

    postBackHome: false,
    addEventToCommentForm: false,
};

export { globals, eventListeners };