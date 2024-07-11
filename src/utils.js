import { allUsers } from "./stores";
import { chatTabs } from "./stores";
import { get } from "svelte/store"


//backend genereerib uuid ja front end paneb clienti session cookie paika.
import Message from './components/chat/message.svelte';
export function updateSessionToken(token, expire) {
    var dateTime = new Date();
    dateTime.setTime(dateTime.getTime() + expire * 60 * 60 * 1000);
    var expires = "expires=" + dateTime.toUTCString();
    document.cookie = "sessionToken=" + token.String + ";" + expires;
  }

export const fetchUsers = async () => {
    const response = await fetch('http://localhost:8080/allusers');
    if(response.ok) {
        const fetchedUsers = await response.json();
        allUsers.set([...fetchedUsers])
        console.log(allUsers)
    } else {
        console.error('Error fetching users:', response.status);
    }
};

export function InsertNewMessage(msgObj) {

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
  //console.log("INSERTIME", msgObj)

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
  PrivateMessageNotification(msgObj.fromUserID)
  if (scrollIsBottom(chatBody, 80)) {
    scrollToBottom(chatBody, false);
  }
}

// ||> PrivateMessage Notification (userlist/open & minimized chats)
function PrivateMessageNotification(fromUserID) {
  const chatContainer = document.getElementById('bottomChatContainer')
  if (!chatContainer) {
      console.error("Couldn't getElementById: #bottomChatContainer")
      return
  }
  const chatBody = chatContainer.querySelector(`div[userid="${fromUserID}"]`)
  if (!chatBody) {
    // IF chat isn't open add a notification to this userID on userlist
    const usersContainer = document.getElementById('usersContainer')
    const targetUserDiv = usersContainer.querySelector(`div[userid="${fromUserID}"]`)
    
    // Edge case for self
    if (!targetUserDiv) {
      return
    }

    targetUserDiv.classList.add('notification')
    const messageIcon = targetUserDiv.querySelector('.messageNotification');
    messageIcon.style.visibility = 'visible';
    return
  }

  const chatPreview = chatBody.querySelector('.chat-preview')
  const previewVisibility = (window.getComputedStyle(chatPreview)).visibility;
  if (previewVisibility === "visible") {
    chatPreview.classList.add('notification');
    return
  } else if (previewVisibility === "collapse") {
    const collapsedBody = chatBody.querySelector('.chat-body')
    const notification = chatBody.querySelector('.new-message-notification')
    if (!scrollIsBottom(collapsedBody, 80)) {
      notification.style.display = 'block';
    }
  } else {
    console.log("Error unexpected value for ChatModule Preview visibility!:", previewVisibility)
  }

}

function scrollToBottom(bodyElem, animation = true) {
  let startTime;
  let start = bodyElem.scrollTop;
  let end = bodyElem.scrollHeight - bodyElem.clientHeight;
  if (animation === false) {
      bodyElem.scrollTop = bodyElem.scrollHeight - bodyElem.clientHeight;
      return
  }
  let duration = 250
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
  return bodyElem.scrollTop >= (bodyElem.scrollHeight - bodyElem.clientHeight - buffer);
}

export function removeFromActiveChat(event, modi='',userID ) {
  // event.stopPropagation();
  // let containerElem = event.target.closest('.chatBox');
  let containerElem = document.querySelector(`.chatBox[userid="${userID}"]`);
  
  // Minimize animation before closing
  let chatPopup = containerElem.querySelector('.chat-popup');
  chatPopup.classList.remove('chat-popup-open')
  chatPopup.classList.add('chat-popup-close')
  // console.log("Removing from active chat");

  if (modi === 'instant') {
      containerElem.classList.add('user-active-chat-remove')
      setTimeout(() => {
          if (containerElem) {
              containerElem.remove();
              chatTabs.update(tabs => tabs.filter(tab => tab.userID !== userID));
              //console.log('chatTabs:', $chatTabs)
          }
      },250)
  }else if (modi === 'openChat') {
    containerElem.classList.add('user-active-chat-remove')
    containerElem.remove();
  } else {
      const chatPreview = containerElem.querySelector('.chat-preview')
      chatPreview.style.visibility = 'visible';
      setTimeout(() => {
          chatPopup.style.display = 'none';
          chatPopup.classList.remove('chat-popup-close');
          containerElem.classList.add('user-active-chat-remove')
          setTimeout(() => {
              if (containerElem) {
              containerElem.remove();
              chatTabs.update(tabs => tabs.filter(tab => tab.userID !== userID));
              console.log('chatTabs:', $chatTabs)
              }
          },220)
      },250)
  }
}

export const getPosts = async () => {
  try {
      const response = await fetch('http://localhost:8080/posts');
      if (response.ok) {
          const fetchedPosts = await response.json();
          posts.set(fetchedPosts); // Update the writable store
          console.log(fetchedPosts)
      } else {
          console.error('Error fetching posts:', response.status);
      }
  } catch (error) {
      console.error('Error:', error);
  }
};

export function getUserDetails(userID) {
  const users = get(allUsers);
  return users.find(user => user.ID === userID);
}

export async function selectUser(userID) {
  const response = await fetch("http://localhost:8080/user?id=" + userID);
  if (response.ok) {
    const selectedUser = await response.json();
    userProfileData.set(selectedUser);
  } else {
    console.error("Error fetching users:", response.status);
  }
}