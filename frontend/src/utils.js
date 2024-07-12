import { allUsers, currentPosts, userProfileData, allGroups, groupSelected, activeTab,events, chatTabs, API_URL, tabMap} from "./stores";
import { get } from 'svelte/store';
import { notifications } from "./websocket.js"

//backend genereerib uuid ja front end paneb clienti session cookie paika.
import Message from "./components/chat/message.svelte";
export function updateSessionToken(token, expire) {
    var dateTime = new Date();
    dateTime.setTime(dateTime.getTime() + expire * 60 * 60 * 1000);
    var expires = "expires=" + dateTime.toUTCString();
    document.cookie = "sessionToken=" + token.String + ";" + expires;
}

export const fetchUsers = async () => {
    const response = await fetch(`${API_URL}/allusers`,{
      credentials: 'include'
    });
    if (response.ok) {
        const fetchedUsers = await response.json();
        allUsers.set([...fetchedUsers])
    } else {
        console.error("Error fetching users:", response.status);
    }
};


export const fetchNotifications = async () => {
    const response = await fetch(`${API_URL}/notifications`,{
      credentials: 'include'
    });
    if (response.ok) {
        const fetchedNotifications = await response.json();
        console.log('alloo')
        console.log(fetchedNotifications.notifications)
        if (fetchedNotifications.notifications === undefined) {
          notifications.update((n) => [...n, ...fetchedNotifications.notifications]);
        }

    } else {
        console.error("Error fetching users: ", response.status);
    }
}

export function InsertNewMessage(msgObj, isGroup) {
  const chatContainer = document.getElementById('bottomChatContainer')
  if (!chatContainer) {
      console.error("Couldn't getElementById: #bottomChatContainer")
      return
  }
  const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)

  switch (isGroup) {
    case true:
        if (!chatBody) {
                // console.error("Got a message, but user hasn't opened this chat, yet, add a notification instead")
                setTimeout(() => {
                    GroupMessageNotification(msgObj.chatID)
                }, 500)
                // PrivateMessageNotification(msgObj.fromUserID)
                return
            }
            // Create the chatBox module
            const GroupMessageElem = new Message({
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
            GroupMessageNotification(msgObj.chatID);
            if (scrollIsBottom(chatBody, 80)) {
                scrollToBottom(chatBody, false);
            }
        
        break;
        default:
            // const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)
            if (!chatBody) {
            // console.error("Got a message, but user hasn't opened this chat, yet, add a notification instead")
            setTimeout(() => {
                PrivateMessageNotification(msgObj.fromUserID)
            }, 500)
            // PrivateMessageNotification(msgObj.fromUserID)
            return
            }
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
  
//   const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)
//   if (!chatBody) {
//     // console.error("Got a message, but user hasn't opened this chat, yet, add a notification instead")
//     setTimeout(() => {
//       PrivateMessageNotification(msgObj.fromUserID)
//     }, 500)
//     // PrivateMessageNotification(msgObj.fromUserID)
//     return
//   }

//   // console.log(msgObj.fromUser)
//   // console.log("here")
//   // console.log(msgObj)

//   // Create the chatBox module
//   const messageElem = new Message({
//       target: chatBody,
//       props: {
//           fromUser: msgObj.fromUserID,
//           fromUsername: msgObj.fromUsername,
//           time: msgObj.time, 
//           msgID: msgObj.msgID,
//           msgContent: msgObj.content,
//           AvatarPath: msgObj.AvatarPath
//       }
//   });

//     // Scrolling and notif logic
//     // This is to prevent instant scroll to bottom when user is mid-scrolling and gets a new message
//     PrivateMessageNotification(msgObj.fromUserID);
//     if (scrollIsBottom(chatBody, 80)) {
//         scrollToBottom(chatBody, false);
//     }
}

// ||> GroupMessage Notification (grouplist/open & minimized chats)
function GroupMessageNotification(chatID) {
    const chatContainer = document.getElementById("bottomChatContainer");
    if (!chatContainer) {
        console.error("Couldn't getElementById: #bottomChatContainer");
        return;
    }
    const chatBody = chatContainer.querySelector(`div[userid="GroupChatID_${chatID}"]`);
    if (!chatBody) { //SSSSSSSSSSSSSSSSSSSSSSSSSs
        // IF chat isn't open add a notification to this chatID on grouplist
        const groupsContainer = document.getElementById("groupsContainer");
        const targetUserDiv = groupsContainer.querySelector(
            `div[groupchatid="${chatID}"]`
        );

        if (!targetUserDiv) {
            return;
        }

        targetUserDiv.classList.add("notification");
        const messageIcon = targetUserDiv.querySelector(".messageNotification");
        messageIcon.style.visibility = "visible";
        return;
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

export function removeFromActiveChat(event, modi='',userID, isGroup ) {
  // event.stopPropagation();
  // let containerElem = event.target.closest('.chatBox');
  let containerElem;
  containerElem = document.querySelector(`.chatBox[userid="${userID}"]`);
/*   if (!isGroup) {
  } else {
    containerElem = document.querySelector(`.chatBox[userid="${userID}"]`);
  } */
  console.log(userID)
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
    console.log("yo",get(tabMap))
    tabMap.update(map => {// SEDA ON VAJA et eemaldada just kustutatud chat tabMapist, kuna tabMap blokeerib chattide ehitamist
      map.delete(userID);  // Perform the deletion
      return map;          // Return the updated map
  });
    
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
      const groupID = get(groupSelected);
      const response = await fetch(`${API_URL}/posts`, {
          credentials: 'include',
          method: 'POST',
          headers: {
              'Content-Type': 'application/json'
          },
          body: groupID,
      });
      if (response.ok) {
          const fetchedPosts = await response.json();
          currentPosts.set(fetchedPosts); // Update the writable store
      } else {
          console.error('Error fetching posts:', response.status);
      }
  } catch (error) {
      console.error('Error:', error);
  }
};
export const getComments = async (postID) => {
  try {
      const response = await fetch(`${API_URL}/comment?postID=${postID}`,{
        credentials: 'include'
      });
      if (response.ok) {
        let comments = await response.json()
        console.log(comments)
        return comments
      }
  } catch (error) {
      console.error('Error:', error);
  }
};

export const getGroups = async () => {
  try {
      const response = await fetch(`${API_URL}/groups`,{
        credentials: 'include'
      });
      if (response.ok) {
          const fetchedGroups = await response.json();
          allGroups.set(fetchedGroups); // Update the writable store
          console.log("GROUP FETCH OK!", fetchedGroups)
      } else {
          console.error('Error fetching posts:', response.status);
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
  const response = await fetch(`${API_URL}/user?id=${userID}`,{
    credentials: 'include'
  });
  if (response.ok) {
    const selectedUser = await response.json();
    activeTab.set('Profile')
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

// Analog from profile.svelte
export async function sendFollow(action, target) {
    try {
        const response = await fetch(`${API_URL}/followers`, {
            credentials: 'include',
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ action: action, target: target }),
        });

        if (!response.ok) {
            throw new Error("Failed to send follow request");
        }
    } catch (error) {
        console.error("Error sending follow request:", error.message);
    }
}



export function leaveGroup(groupID){
fetch(`${API_URL}/leaveGroup`, {
  credentials: 'include',
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    groupID: groupID,
  },)
}).then((response) => {
  if (response.ok) {
    console.log("Group left");
    getGroups();
    groupSelected.set(0)
    groupSelected.set(groupID)
  }
}).catch((error) => {
  console.error("Error leaving group:", error);
});
}

export const joinGroup = (groupID, action) => {
  fetch(`${API_URL}/joinGroup`, {
    credentials: 'include',
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      groupID: groupID,
      action: action,
    }),
  })
    .then((response) => {
      if (response.ok) {
        response.json().then((data) => {
          console.log(data);
          getGroups();
          groupSelected.set(0); // to force reactivity
          groupSelected.set(data.groupID);
        });
      }
    })
    .catch((error) => {
      console.error("Error joining group:", error);
    });
};

export async function getEvents(groupID) {
    try {
      const response = await fetch(`${API_URL}/events`, {
        credentials: 'include',
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: groupID,
      });
      if (response.ok) {
        const fetchedEvents = await response.json();
        events.set(fetchedEvents);
      } else {
        console.error("Error fetching events:", response.status);
      }
    } catch (error) {
      console.error("Error:", error);
    }
  }

  export   function deleteGroup(groupID) {
    fetch(`${API_URL}/deleteGroup`, {
      credentials: 'include',
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        groupID: groupID,
      }),
    })
      .then((data) => {
        console.log(data);
        getGroups();
        groupSelected.set(0);
        getPosts();
      })
      .catch((error) => console.error(error));
  }