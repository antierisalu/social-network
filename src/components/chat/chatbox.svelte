<script>
  // variables to create the chatbox
  let isTyping = true;
  const userID = 0;
  const chatID = 0;

  const userName = "Pepe Frog";
  // import svg elements
  import CloseChat from "../icons/closeChat.svelte";
  import MinimizeChat from "../icons/minimizeChat.svelte";
  import ChatModuleEmojiPicker from "../icons/chatModuleEmojiPicker.svelte";
</script>

<div
  class="chatBox"
  userid="0"
  id="activeChat-chatModule"
  style="display: flex;"
>
  <div class="chat-popup chat-popup-open">
    <div class="chat-header">
      <div class="wrapper">
        <div class="statusBubble">
          <div class="isTyping">
            <a>is typing</a>
            <div class="typingAnimation">
              <div class="circle c01"></div>
              <div class="circle c02"></div>
              <div class="circle c03"></div>
            </div>
          </div>
        </div>
        <div class="username">
          <a>JamesJackson34</a>
        </div>
      </div>
      <div class="btn-wrapper">
        <!-- Hide/Minimize current chat -->
        <div class="minimize-chat">
          <MinimizeChat />
        </div>
        <!-- Close/Remove current chat -->
        <div class="close-chat">
          <CloseChat />
        </div>
      </div>
    </div>
    <div class="chat-body" chatID="" messageCount=""></div>
    <div class="chat-footer">
      <div contenteditable class="chatModule-input-field"></div>
      <div class="chatModule-emoji-picker">
        <ChatModuleEmojiPicker />
        <!--   <svg width="32px" height="32px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" transform="rotate(0)matrix(1, 0, 0, 1, 0, 0)">
                    <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                    <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
                    <g id="SVGRepo_iconCarrier">
                        <path opacity="0.5" d="M6.50359 21.5201C8.13698 22.4631 10.2103 21.9075 14.357 20.7964C18.5037 19.6853 20.577 19.1298 21.5201 17.4964C22.4631 15.863 21.9075 13.7897 20.7964 9.643C19.6853 5.49632 19.1298 3.42298 17.4964 2.47995C15.863 1.53691 13.7897 2.09246 9.643 3.20356C5.49632 4.31466 3.42298 4.87021 2.47995 6.50359C1.53691 8.13698 2.09246 10.2103 3.20356 14.357C4.31466 18.5037 4.87021 20.577 6.50359 21.5201Z" fill="#adff2f"></path>
                        <path d="M14.8978 11.2237C15.4313 11.0808 15.6899 10.3162 15.4755 9.51599C15.2611 8.71579 14.6548 8.18298 14.1213 8.32592C13.5879 8.46886 13.3292 9.23343 13.5436 10.0336C13.7581 10.8338 14.3643 11.3666 14.8978 11.2237Z" fill="#adff2f"></path>
                        <path d="M9.10238 12.7767C9.63585 12.6337 9.89449 11.8692 9.68008 11.069C9.46567 10.2688 8.85939 9.73596 8.32592 9.8789C7.79246 10.0218 7.53381 10.7864 7.74823 11.5866C7.96264 12.3868 8.56892 12.9196 9.10238 12.7767Z" fill="#adff2f"></path>
                        <path d="M8.18524 15.751C8.28594 15.3492 8.69329 15.1052 9.09507 15.2059C10.2254 15.4892 11.5234 15.4927 12.8411 15.1396C14.1589 14.7865 15.2813 14.1345 16.1185 13.324C16.4161 13.0359 16.8909 13.0436 17.179 13.3412C17.4671 13.6388 17.4594 14.1136 17.1618 14.4017C16.8143 14.7381 16.4298 15.0495 16.0129 15.3304L16.1709 15.6523C16.5396 16.4034 16.2225 17.3108 15.4663 17.6688C14.7251 18.0197 13.8395 17.7102 13.4781 16.9741L13.2819 16.5742L13.2294 16.5885C11.674 17.0052 10.1168 17.0083 8.73039 16.6609C8.32861 16.5602 8.08453 16.1528 8.18524 15.751Z" fill="#adff2f"></path>
                    </g>
                </svg> -->
      </div>
      <!-- <div class="chatModule-input-send"> 
                <svg width="38px" height="38px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                    <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
                    <g id="SVGRepo_iconCarrier"> 
                        <path d="M11.5003 12H5.41872M5.24634 12.7972L4.24158 15.7986C3.69128 17.4424 3.41613 18.2643 3.61359 18.7704C3.78506 19.21 4.15335 19.5432 4.6078 19.6701C5.13111 19.8161 5.92151 19.4604 7.50231 18.7491L17.6367 14.1886C19.1797 13.4942 19.9512 13.1471 20.1896 12.6648C20.3968 12.2458 20.3968 11.7541 20.1896 11.3351C19.9512 10.8529 19.1797 10.5057 17.6367 9.81135L7.48483 5.24303C5.90879 4.53382 5.12078 4.17921 4.59799 4.32468C4.14397 4.45101 3.77572 4.78336 3.60365 5.22209C3.40551 5.72728 3.67772 6.54741 4.22215 8.18767L5.24829 11.2793C5.34179 11.561 5.38855 11.7019 5.407 11.8459C5.42338 11.9738 5.42321 12.1032 5.40651 12.231C5.38768 12.375 5.34057 12.5157 5.24634 12.7972Z" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="--darkreader-inline-stroke: #e8e6e3;" data-darkreader-inline-stroke=""></path>
                    </g>
                </svg>
            </div> -->
    </div>
    <!-- <div class="new-message-notification2 typingGlow" style="display:none">
        </div>
        <div class="new-message-notification" style="display: none;">
            <div class="notif-wrapper" onclick="scrollChatBottom(event)">
                <svg width="31px" height="31px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" stroke="#ffffff" style="--darkreader-inline-stroke: #e8e6e3;" data-darkreader-inline-stroke="">
                    <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                    <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
                    <g id="SVGRepo_iconCarrier"> 
                        <path d="M9 13L12 16M12 16L15 13M12 16V8M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="--darkreader-inline-stroke: #e8e6e3;" data-darkreader-inline-stroke=""></path>
                    </g>
                </svg>
            </div>
        </div> -->
  </div>
  <div class="chat-preview" onclick="toggleChat(event)">
    <div class="wrapper">
      <div class="statusBubble"></div>
      <div class="username">
        <a id="preview-username">JamesJackson34</a>
      </div>
    </div>
    <div class="btn-wrapper">
      <div class="close-chat">
        <CloseChat />
      </div>
    </div>
  </div>
  <script>
    function minimizeChat(event) {
      const chatPopup = event.closest(".chat-popup");
      const chatPreview = chatPopup.nextElementSibling;
      chatPopup.classList.remove("chat-popup-open");
      chatPopup.classList.add("chat-popup-close");
      setTimeout(() => {
        chatPopup.style.display = "none";
        chatPopup.classList.remove("chat-popup-close");
        chatPopup.classList.add("chat-popup-open");
      }, 310);
      chatPreview.style.visibility = "visible";
    }
    function toggleChat(event) {
      const chatPreview = event.currentTarget.closest(".chat-preview");
      chatPreview.style.visibility = "collapse";
      const chatPopup = chatPreview.previousElementSibling;
      chatPopup.style.display = "flex";
      // Notification remove on-click
      chatPreview.classList.remove("notification");
      // Check if its the first load (SCROLLING)
      const activeChat = event.currentTarget.closest(".chatBox");
      if (activeChat && activeChat.hasAttribute("firstload")) {
        activeChat.removeAttribute("firstload");
        const chatBody = chatPopup.querySelector(".chat-body");
        chatBody.scrollTop = chatBody.scrollHeight - chatBody.clientHeight;
      }
    }
    // Scrolls to bottom with/without animation
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
    function scrollChatBottom(event) {
      const activeChat = event.currentTarget.closest(".chatBox");
      const chatBody = activeChat.querySelector(".chat-body");
      scrollToBottom(chatBody);
      const notification = activeChat.querySelector(
        ".new-message-notification"
      );
      notification.style.display = "none";
    }
    function removeFromActiveChat(event, modi = "") {
      event.stopPropagation();
      let containerElem = event.target.closest(".chatBox");
      // Minimize animation before closing
      let chatPopup = containerElem.querySelector(".chat-popup");
      chatPopup.classList.remove("chat-popup-open");
      chatPopup.classList.add("chat-popup-close");
      // console.log("Removing from active chat");
      if (modi === "instant") {
        containerElem.classList.add("user-active-chat-remove");
        setTimeout(() => {
          if (containerElem) {
            containerElem.remove();
          }
        }, 250);
      } else {
        const chatPreview = containerElem.querySelector(".chat-preview");
        chatPreview.style.visibility = "visible";
        setTimeout(() => {
          chatPopup.style.display = "none";
          chatPopup.classList.remove("chat-popup-close");
          containerElem.classList.add("user-active-chat-remove");
          setTimeout(() => {
            if (containerElem) {
              containerElem.remove();
            }
          }, 220);
        }, 250);
      }
    }
  </script>
</div>

<style>
  /* ajutine: */
  #preview-username,
  .username a {
    color: #ffffff;
    text-decoration: none;
    font-size: 16px;
  }
  :root {
    --chatWidth: 264px;
    --chatPreviewH: 40px;
    --chatFullH: 400px;
  }

  .chatBox {
    margin-right: 6px;
    margin-left: 6px;
    width: var(--chatWidth);
    height: (--chatPreviewH);
    display: flex;
    flex-direction: row;
    margin-right: 5px;
  }

  /* unused */
  /*  .activeChats {
    display: flex;
    justify-content: end;
    align-items: center;
    flex-direction: row-reverse;
    position: absolute;
    bottom: 10px;
    left: 0;
    width: 100%;
    height: 72px;
    } */

  /* unused */
  /*   .user-active-chat {
        display: flex;
        flex-direction: row;
        margin: 5px;
        margin-right: 10px;
        width: 316px;
        height: 50px;
    } */
  .chat-preview {
    display: flex;
    justify-content: space-between;
    width: var(--chatWidth);
    height: var(--chatPreviewH);
    background-color: black;
    border-radius: 5px;
    border: 1px solid rgb(145, 145, 145);
  }
  .chat-preview .wrapper {
    width: fit-content;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .chat-preview .btn-wrapper {
    height: 100%;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .chat-popup {
    display: none;
    flex-direction: column;
    position: absolute;
    transform: translatey(-350px);
    width: var(--chatWidth);
    height: var(--chatFullH);
    background-color: black;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.125);
  }
  .chat-popup-open {
    animation: chat-open 300ms ease-in-out forwards;
    transform-origin: bottom;
  }

  :global(.chat-popup-close) {
    animation: chat-close 300ms ease-in-out forwards;
    transform-origin: bottom;
  }
  .chatModule-emoji-picker {
    margin-right: 3px;
    width: 34px;
    height: 34px;
  }

  .statusBubble {
    margin-left: 15px;
    margin-right: 10px;
    background-color: #168feb;
    width: 18px;
    height: 18px;
    border-radius: 100%;
    border: 3px solid #1584d8e0;
    -webkit-box-shadow: 0px 0px 10px 2px rgba(1, 102, 254, 0.49);
    -moz-box-shadow: 0px 0px 10px 2px rgba(1, 102, 254, 0.49);
    box-shadow: 0px 0px 10px 2px rgba(1, 102, 254, 0.49);
    transition:
      background-color 0.2s ease,
      border-color 0.2s ease,
      box-shadow 0.2s ease;
  }
  .chat-header {
    display: flex;
    justify-content: space-between;
    width: var(--chatWidth);
    height: 60px;
    background-color: rgba(17, 25, 40, 0.75);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.125);
    border-bottom-left-radius: 1px;
    border-bottom-right-radius: 1px;
  }

  .chat-header .wrapper {
    width: fit-content;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 12px;
  }

  .chat-header .btn-wrapper {
    height: 100%;
    position: relative;
    right: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .chat-header .statusBubble {
    margin-left: 0;
  }

  .chat-body {
    width: var(--chatWidth);
    height: 280px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    overflow: hidden scroll;
  }
  :global(.chat-time-wrapper) {
    width: 100%;
    display: flex;
    justify-content: end;
    margin-top: 2px;
  }
  :global(.message-container) {
    margin: 10px;
    margin-top: 2px;
    margin-bottom: 2px;
    width: 290px;
    min-height: var(--chatPreviewH);
    height: fit-content;
  }
  :global(.message-container:hover .chat-time) {
    top: 0;
    opacity: 1;
  }
  :global(.message-header) {
    display: flex;
    justify-content: space-between;
  }
  :global(.chat-time) {
    position: relative;
    color: gray;
    font-weight: 500;
    min-height: 18px;
    user-select: none;
    opacity: 0;
    transition:
      top 0.3s ease,
      opacity 0.3s ease;
  }
  :global(.chat-username) {
    font-size: large;
    color: white;
    font-weight: 700;
    min-height: 18px;
    user-select: none;
  }
  :global(.chat-message-content) {
    font-size: medium;
    color: white;
    font-weight: 600;
  }
  .chat-footer {
    position: absolute;
    bottom: 0;
    left: 0;
    width: var(--chatWidth);
    height: var(--chatPreviewH);
    background-color: rgba(17, 25, 40, 0.95);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.125);
    display: flex;
    flex-direction: row;
    align-items: center;
    width: 100%;
    min-height: var(--chatPreviewH);
    height: fit-content;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
    z-index: 2;
  }
  .chatModule-input-field {
    margin: 5px;
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    background-color: rgba(12, 18, 31, 0.75);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.125);
    width: 220px;
    text-align: left;
  }
  .chatModule-input-field[contenteditable] {
    color: white;
    font-weight: 1000;
    min-height: 18px;
    padding: 8px;
    padding-top: 4px;
    padding-bottom: 6px;
    line-height: 1.4;
    border-radius: 10px;
  }
  .chatModule-input-field[contenteditable]:focus {
    outline: none;
  }
  .chatModule-input-send {
    position: absolute;
    bottom: 0;
    right: 0;
    margin: 6px;
    margin-right: 15px;
    width: fit-content;
    height: fit-content;
  }

  @keyframes user-active-chat-remove {
    0% {
      opacity: 1;
      transform: scaleY(1);
      transform: translatey(0px);
    }
    100% {
      opacity: 0;
      transform: scaleY(0);
      transform: translatey(-20px);
    }
  }
  :global(.user-active-chat-remove) {
    animation: user-active-chat-remove 300ms ease-in-out forwards;
    transform-origin: top;
  }
  @keyframes chat-open {
    0% {
      opacity: 0;
      height: var(--chatPreviewH);
      transform: scaleY(0);
      transform: translatey(0px);
    }
    100% {
      opacity: 1;
      height: var(--chatFullH);
      transform: scaleY(1)
        translateY(calc(-1 * (var(--chatFullH) - var(--chatPreviewH))));
    }
  }

  @keyframes chat-close {
    0% {
      opacity: 1;
      height: var(--chatFullH);
      transform: scaleY(1)
        translateY(calc(-1 * (var(--chatFullH) - var(--chatPreviewH))));
    }
    100% {
      opacity: 0;
      height: var(--chatPreviewH);
      transform: scaleY(0);
      transform: translatey(0px);
    }
  }

  :global(.full-name) {
    position: absolute;
    transform: translate(20px, 14px);
    font-family: "Jura";
    color: white;
    font-size: small;
    font-weight: 1000;
  }

  .isTyping {
    display: none;
    text-wrap: nowrap;
    align-items: center;
    justify-content: start;
    width: 200px;
    position: relative;
    transform: translate(20px, 14px);
    font-family: "Jura-Bold";
    color: white;
    font-size: small;
    font-weight: 700;
  }

  .typingAnimation {
    transform: translate(-1px, -1px);
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
    align-items: end;
    width: 19px;
    height: 12px;
  }
  .circle {
    background: white;
    width: 2.5px;
    height: 2.5px;
    border-radius: 100%;
    animation: wave 0.9s infinite;
  }
  .c01 {
    animation-delay: 0.1s;
  }
  .c02 {
    animation-delay: 0.2s;
  }
  .c03 {
    animation-delay: 0.3s;
  }
  @keyframes wave {
    0% {
      transform: translateY(-0.5px);
    }
    45% {
      transform: translateY(-4px);
    }
    100% {
      transform: translateY(0px);
    }
  }
  :global(.typingGlow) {
    display: none;
    background: linear-gradient(
      0deg,
      rgba(150, 4, 254, 0.645) 43%,
      rgba(178, 4, 254, 0) 92%
    );
    animation: pulseGlow 1.5s infinite;
  }
  @keyframes pulseGlow {
    0% {
      opacity: 0.2;
    }
    50% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
</style>
