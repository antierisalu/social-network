//when the value of a store changes, all components
// that are subscribed to that store will be changed
import { writable} from 'svelte/store'
import { sendMessage } from './websocket';

export const API_URL = "http://localhost:8080/api"
export const IMAGE_URL = "http://localhost:8080"

export const groupSelected = writable(0)

//user is not logged in 
export const loggedIn = writable(false);

// Left sidebar tabs
export const activeTab = writable('Profile')

// Profile searchbar
export const userProfileData = writable({})

// Profile editing
export const isEditingProfile = writable(false);

export const newAboutMeStore = writable('')

export const allUsers = writable([])

//this is client's own info
export const userInfo = writable({})

// Auth errors (login/register)
export const authError = writable('');
export function displayUserAuthError(errorStr) {
    authError.set(errorStr)
    setTimeout(() => {
      authError.set('')
    }, 3000);
}

//Posts of last selected group or main feed.
export const currentPosts = writable([]);

export const allGroups = writable([]);
window.allGroups = allGroups;

export const events = writable([])

export const uploadImageStore = writable(null)

// Connected with WS Online users (user ID's)
export const onlineUserStore = writable([])

// Connected with WS (last messages with target userID)
// map[userID][lastMessage]
export const lastMsgStore = writable({})

// Contains chat notification states (seen)
// map[userID (int)][lastUnseenMessage(int)]
export const chatNotifStore = writable({})
window.chatNotifStore = chatNotifStore;

// Frontend store update + backend 
// Mark current message and all of the prior messages to seen for Client(store) + update DB (seen)
// Note: userID is targetID (not chatID, chatID is derived from userID and fromID on backend)
export function markMessageAsSeen(userID) {
  // Remove the chatNotifstore entry that has the k,v pair with value of userID
  let messageID;
  let fromID;
  // Update store
  chatNotifStore.update(store => {
    const { [userID]: val, ...newStore} = store;
    messageID = val;
    return newStore;
  });
  userInfo.subscribe(userInfo => {
    fromID = userInfo.id
  });
  // Backend. This causes circular dependency.
  sendMessage(JSON.stringify({ type: "markAsSeen", id: userID, targetID: messageID, fromID: fromID }))
}

// Contains Group chat ID's that have an un-resolved notification (seen)
export const groupChatNotifStore = writable([])

export function markGroupMessageAsSeen(chatID) {
  let fromID;
  let group;
  userInfo.subscribe(userInfo => {
    fromID = userInfo.id
  });
  allGroups.subscribe(groupstore => {
    group = groupstore.find(group => group.chatid === chatID)
  });

  // Update store
  groupChatNotifStore.update(store => {
    const newStore = store.filter(id => id !== chatID);
    console.log("Store updated!")
    console.log(newStore)
    console.log(fromID)
    console.log("GroupID:", group.id)

    return newStore;
  });

  



  sendMessage(JSON.stringify({ type: "markGroupAsSeen", targetID: group.id, fromID: fromID }))
  
}

// Contains all the users currently typing to client 
export const isTypingStore = writable([])
const typingTimeouts = new Map();
// Utility function to add value & remove it after 2 sec
export function setTyping(userID) {
  if (typingTimeouts.has(userID)) {
    clearTimeout(typingTimeouts.get(userID));
  }
  isTypingStore.update(current => {
    const idx = current.findIndex(item => item === userID);
    if (idx !== -1) {
      current.splice(idx, 1);
    }
    current.push(userID)
    return current;
  });
  const timeoutID = setTimeout(() => {
    removeTyping(userID)

  }, 2000);
  typingTimeouts.set(userID, timeoutID)
}

export function removeTyping(userID) {
  isTypingStore.update(current => {
    const idx = current.findIndex(item => item === userID);
    if (idx !== -1) {
      current.splice(idx, 1);
    }
    return current;
  });
  typingTimeouts.delete(userID)
}

//Array to store the opened chatbox tabs.
export const chatTabs = writable([]);
