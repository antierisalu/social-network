//when the value of a store changes, all components
// that are subscribed to that store will be changed
import { writable } from 'svelte/store'
import { sendMessage } from './websocket';

//user is not logged in 
export const loggedIn = writable(false);

export const activeTab = writable('Profile')

// Profile searchbar
export const userProfileData = writable({})

// Profile editing
export const isEditingProfile = writable(false);

export const newAboutMeStore = writable('')
export const showOverlay = writable(false)

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

export const uploadImageStore = writable(null)

// Connected with WS Online users (user ID's)
export const onlineUserStore = writable([])

// Connected with WS (last messages with target userID)
// map[userID][lastMessage]
export const lastMsgStore = writable({})

// Contains chat notification states (seen)
// map[userID (int)][lastUnseenMessage(int)]
export const chatNotifStore = writable({})

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

  // Backend
  sendMessage(JSON.stringify({ type: "markAsSeen", id: userID, targetID: messageID, fromID: fromID }))
}