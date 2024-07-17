<script>
    import MsgNotification from "../icons/msgNotification.svelte";
    import { connect, sendMessage, messages, sendDataRequest } from "../../websocket";
    import { get } from "svelte/store";
    import { activeTab, isTypingStore, userInfo } from "../../stores";
    import Message from './message.svelte';
    import Chatbox from "./chatbox.svelte";
    import { allUsers, markMessageAsSeen,IMAGE_URL, allowedTabAmount } from "../../stores";
    import { chatTabs } from "../../stores";
    import { removeFromActiveChat } from "../../utils";


    $: users = $allUsers;
    export let avatarPath = "";
    if (avatarPath === "") {
        avatarPath = "/images/avatars/default.png"
    }
    export let firstName = "";
    export let lastName = "";
    export let userID = "";
    export let isOnline;
    export let lastNotification;
    let chatID;
    $: typingStore = $isTypingStore
    
    function removeNotificationClass(userID) {
        const usersContainer = document.getElementById('usersContainer')
        const targetUserDiv = usersContainer.querySelector(`div[userid="${userID}"]`)
        targetUserDiv.classList.add('notification')
        if (targetUserDiv) {
            targetUserDiv.classList.remove('notification')
            const messageIcon = targetUserDiv.querySelector('.messageNotification');
            messageIcon.style.visibility = 'hidden';
        }

        // [Frontend + Backend] Remove from chatNotifStore (userID) && send through WS (to mark all messages to seen to last notif message)
        markMessageAsSeen(userID)
        // ^ This can be added to bottom chat-modules later on as needed, currently just for the allUsers tab.
    }

    export function addToChatTabsArray(userID, firstName, lastName, avatarPath) {
        const existTab = $chatTabs.some(tab => tab.userID === userID);

        if (!existTab) {
            $chatTabs = [...$chatTabs, { userID, firstName, lastName, avatarPath, isOnline }];
            
            if ($chatTabs.length > $allowedTabAmount) {
                const removedUserID = $chatTabs[$chatTabs.length-3].userID    
                removeFromActiveChat(event, 'openChat', removedUserID);
            }
        } else {
            console.log(`userID already exist in chatTab array.`);
        }
    }

    function handleClick() {
        addToChatTabsArray(userID, firstName, lastName, avatarPath);
        removeNotificationClass(userID);
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="user {(typeof lastNotification === "number") ? 'notification' : ''}" {userID} on:click={handleClick}>
    <div class="profilePictureWrapper  {(isOnline) ? 'online' : 'offline'}">
        <img src={IMAGE_URL}{avatarPath} alt={userID} class="{(isOnline) ? '' : 'avatar-grayscale'}">
    </div>

    <div class="usernameWrapper">
        <h2 class="username" style="margin: 0;">{firstName} {lastName}</h2>
    </div>

    <div class="messageNotification">
        <MsgNotification />
    </div>
</div>

<style>
    .user {
        user-select: none;
        cursor: pointer;
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 3%;
        width: 94%;
        height: 42px;
        border-radius: 5px;
        border: 1px solid rgb(145, 145, 145);
    }

    .profilePictureWrapper {
        margin-left: 6px;
        margin-top: 3px;
        margin-bottom: 3px;
        width: 34px;
        height: 34px;
        border-radius: 50%;
        /* border: 2px solid #5f9313bd; */
    }

    :global(.online),
    :global(.offline) {
        border: 2px solid #636363;
        transition: border-color 0.3s ease;
    }
    :global(.online) {
        border-color: #2ccc00c9;
    }
    :global(.offline) {
        border-color: #636363;
    }
    :global(.avatar-grayscale) {
        filter: grayscale(100%);
        transition: filter 0.3s ease;
    }
    
    .profilePictureWrapper img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 50%;
    }

    .usernameWrapper {
        display: flex;
        justify-content: center;
        align-items: center;
        flex: 7;
        height: 100%;
    }

    .username {
        font-size: medium;
    }

    .messageNotification {
        visibility: hidden;
        flex: 1;
        margin-right: 6px;
        margin-top: 3px;
        margin-bottom: 3px;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
    }
</style>