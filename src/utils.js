import { allUsers, allPosts, userProfileData } from "./stores";
import { get } from "svelte/store";

//backend genereerib uuid ja front end paneb clienti session cookie paika.
import Message from "./components/chat/message.svelte";
export function updateSessionToken(token, expire) {
    var dateTime = new Date();
    dateTime.setTime(dateTime.getTime() + expire * 60 * 60 * 1000);
    var expires = "expires=" + dateTime.toUTCString();
    document.cookie = "sessionToken=" + token.String + ";" + expires;
}

export const fetchUsers = async () => {
    const response = await fetch("http://localhost:8080/allusers");
    if (response.ok) {
        const fetchedUsers = await response.json();
        console.log(fetchedUsers)
        allUsers.set([...fetchedUsers])
    } else {
        console.error("Error fetching users:", response.status);
    }
};

export function InsertNewMessage(msgObj) {
  const chatContainer = document.getElementById('bottomChatContainer')
  if (!chatContainer) {
      console.error("Couldn't getElementById: #bottomChatContainer")
      return
  }
  const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)
  if (!chatBody) {
    // console.error("Got a message, but user hasn't opened this chat, yet, add a notification instead")
    setTimeout(() => {
      PrivateMessageNotification(msgObj.fromUserID)
    }, 500)
    // PrivateMessageNotification(msgObj.fromUserID)
    return
  }

  // console.log(msgObj.fromUser)
  // console.log("here")
  // console.log(msgObj)

  // Create the chatBox module
  const messageElem = new Message({
      target: chatBody,
      props: {
          fromUser: msgObj.fromUserID,
          fromUsername: msgObj.fromUsername,
          time: msgObj.time, 
          msgID: msgObj.msgID,
          msgContent: msgObj.content,
          AvatarPath: msgObj.AvatarPath
      }
  });

    // Scrolling and notif logic
    // This is to prevent instant scroll to bottom when user is mid-scrolling and gets a new message
    PrivateMessageNotification(msgObj.fromUserID);
    if (scrollIsBottom(chatBody, 80)) {
        scrollToBottom(chatBody, false);
    }
}

// ||> PrivateMessage Notification (userlist/open & minimized chats)
function PrivateMessageNotification(fromUserID) {
    const chatContainer = document.getElementById("bottomChatContainer");
    if (!chatContainer) {
        console.error("Couldn't getElementById: #bottomChatContainer");
        return;
    }
    const chatBody = chatContainer.querySelector(`div[userid="${fromUserID}"]`);
    if (!chatBody) {
        // IF chat isn't open add a notification to this userID on userlist
        const usersContainer = document.getElementById("usersContainer");
        const targetUserDiv = usersContainer.querySelector(
            `div[userid="${fromUserID}"]`
        );

        // Edge case for self
        if (!targetUserDiv) {
            return;
        }

        targetUserDiv.classList.add("notification");
        const messageIcon = targetUserDiv.querySelector(".messageNotification");
        messageIcon.style.visibility = "visible";
        return;
    }

    const chatPreview = chatBody.querySelector(".chat-preview");
    const previewVisibility = window.getComputedStyle(chatPreview).visibility;
    if (previewVisibility === "visible") {
        chatPreview.classList.add("notification");
        return;
    } else if (previewVisibility === "collapse") {
        const collapsedBody = chatBody.querySelector(".chat-body");
        const notification = chatBody.querySelector(
            ".new-message-notification"
        );
        if (!scrollIsBottom(collapsedBody, 80)) {
            notification.style.display = "block";
        }
    } else {
        console.log(
            "Error unexpected value for ChatModule Preview visibility!:",
            previewVisibility
        );
    }
}

function scrollToBottom(bodyElem, animation = true) {
    let startTime;
    let start = bodyElem.scrollTop;
    let end = bodyElem.scrollHeight - bodyElem.clientHeight;
    if (animation === false) {
        bodyElem.scrollTop = bodyElem.scrollHeight - bodyElem.clientHeight;
        return;
    }
    let duration = 250;
    function animateScroll(timestamp) {
        if (!startTime) startTime = timestamp;
        const elapsed = timestamp - startTime;
        const progress = Math.min(elapsed / duration, 1);
        bodyElem.scrollTop = start + (end - start) * progress;
        if (elapsed < duration) {
            requestAnimationFrame(animateScroll);
        }
    }
    requestAnimationFrame(animateScroll);
}

// Checks if scroll is at bottom with a buffer
function scrollIsBottom(bodyElem, buffer = 60) {
    return (
        bodyElem.scrollTop >=
        bodyElem.scrollHeight - bodyElem.clientHeight - buffer
    );
}

export const getPosts = async () => {
  try {
      const response = await fetch('http://localhost:8080/posts');
      if (response.ok) {
          const fetchedPosts = await response.json();
          allPosts.set(fetchedPosts); // Update the writable store
      } else {
          console.error('Error fetching posts:', response.status);
      }
  } catch (error) {
      console.error('Error:', error);
  }
};
export const getComments = async (postID) => {
  try {
      const response = await fetch(`http://localhost:8080/comment?postID=${postID}`);
      if (response.ok) {
        let comments = await response.json()
        console.log(comments)
        return comments
      }
  } catch (error) {
      console.error('Error:', error);
  }
};

export function getUserDetails(userID) {
    const users = get(allUsers);
    return users.find((user) => user.ID === userID);
}

export async function selectUser(userID) {
    const response = await fetch("http://localhost:8080/user?id=" + userID);
    if (response.ok) {
        const selectedUser = await response.json();
        console.log(selectedUser)
        userProfileData.set(selectedUser);
    } else {
        console.error("Error fetching users:", response.status);
    }
}

export function bellNotif() {
    const element = document.getElementById("notifbell");
    if (element) {
        element.classList.add("notification");
        console.log("Notification class added to the element");

        element.addEventListener("click", function () {
            this.classList.remove("notification");
            console.log("Notification class removed from the element");
        });
    } else {
        console.error("Element not here");
    }
}
