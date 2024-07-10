<script>
    import MsgNotification from "../icons/msgNotification.svelte";
    import { connect, sendMessage, messages, sendDataRequest } from "../../websocket";
    import { get } from "svelte/store";
    import { activeTab, allGroups, isTypingStore, userInfo } from "../../stores";
    import Message from './message.svelte';
    import Chatbox from "./chatbox.svelte";
    import { allUsers, markMessageAsSeen } from "../../stores";
    import { chatTabs } from "../../stores";
    import { identity } from "svelte/internal";

    // $: users = $allUsers;
    $: groups = $allGroups;

    export let avatarPath = "";
    if (avatarPath === "") {
        avatarPath = "./avatars/defaultGroup.png"
    }
    export let groupTitle = "";


    // export let userID = "";
    export let groupChatID;
    // This is done to avoid clashing with userIDs in chatTabs & to generate UID
    let groupPrefixID = 'GroupChatID_'+groupChatID;


    // export let isOnline;
    export let lastNotification;
    let chatID;
    $: typingStore = $isTypingStore
    
    // function removeNotificationClass(userID) {
    //     const usersContainer = document.getElementById('usersContainer')
    //     const targetUserDiv = usersContainer.querySelector(`div[userid="${userID}"]`)
    //     targetUserDiv.classList.add('notification')
    //     if (targetUserDiv) {
    //         targetUserDiv.classList.remove('notification')
    //         const messageIcon = targetUserDiv.querySelector('.messageNotification');
    //         messageIcon.style.visibility = 'hidden';
    //     }

    //     // [Frontend + Backend] Remove from chatNotifStore (userID) && send through WS (to mark all messages to seen to last notif message)
    //     markMessageAsSeen(userID)
    //     // ^ This can be added to bottom chat-modules later on as needed, currently just for the allUsers tab.
    // }

    export function addToChatTabsArray(userID, firstName, lastName, avatarPath, isGroup, groupChatID) {

        const existTab = $chatTabs.some(tab => tab.userID === userID);

        if (!existTab) {
            $chatTabs = [...$chatTabs, { userID, firstName, lastName, avatarPath, isGroup, groupChatID }];
        }else {
            console.log(`userID already exist in chatTab array.`);
        }
    }

    function handleClick() {

        addToChatTabsArray(groupPrefixID, groupTitle, "", avatarPath, true, groupChatID);
        // removeNotificationClass(groupPrefixID);
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="group {(typeof lastNotification === "number") ? 'notification' : ''}" {groupChatID} on:click={handleClick}>
    <div class="profilePictureWrapper">
        <img src={avatarPath} alt={groupChatID} class="borderColor">
    </div>

    <div class="usernameWrapper">
        <h2 class="username" style="margin: 0;">{groupTitle}</h2>
    </div>

    <div class="messageNotification">
        <MsgNotification />
    </div>
</div>

<style>
    .group {
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
    }


    :global(.borderColor) {
        border: 2px solid;
        border-color: #00cc92c9;
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