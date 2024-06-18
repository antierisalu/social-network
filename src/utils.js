import Message from './components/chat/message.svelte';
export function updateSessionToken(token, expire) {
    var dateTime = new Date();
    dateTime.setTime(dateTime.getTime() + expire * 60 * 60 * 1000);
    var expires = "expires=" + dateTime.toUTCString();
    document.cookie = "sessionToken=" + token.String + ";" + expires;
  }

export function InsertNewMessage(msgObj) {
  console.log("wadap", msgObj)
  const chatContainer = document.getElementById('bottomChatContainer')
  if (!chatContainer) {
      console.error("Couldn't getElementById: #bottomChatContainer")
      return
  }
  const chatBody = chatContainer.querySelector(`div[chatid="${msgObj.chatID}"]`)
  console.log(msgObj.fromUser)
  console.log("here")
  console.log(msgObj)
    // Create the chatBox module
  const messageElem = new Message({
      target: chatBody,
      props: {
          fromUser: msgObj.fromUser,
          fromUsername: msgObj.fromUsername,
          time: msgObj.time,
          msgID: msgObj.msgID,
          msgContent: msgObj.content
      }
  });
}