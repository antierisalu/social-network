import { writable } from 'svelte/store'


export const loggedIn = writable(false);

export const activeTab = writable('Groups')
// export const activeTab = writable('Profile')
// return back to normal(Profile) after dev

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

// User list (right-side bar)
// export const userList = writable([])


// export function InsertNewMessage(msgObj) {
//   console.log("wadap", msgObj)
//   const chatContainer = document.getElementById('bottomChatContainer')
//   if (!chatContainer) {
//       console.error("Couldn't getElementById: #bottomChatContainer")
//       return
//   }
//   const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)
//    // Create the chatBox module
//   const messageElem = new Message({
//       target: chatBody,
//       props: {
//           fromUser: msgObj.fromUser,
//           time: msgObj.time,
//           msgID: msgObj.msgID,
//           msgContent: msgObj.content
//       }
//   });
// }