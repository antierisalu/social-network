<script>
    import MsgNotification from "../icons/msgNotification.svelte";
    import { connect, sendMessage, messages, sendDataRequest } from "../../websocket";
    import { get } from "svelte/store";
    import { activeTab, userInfo } from "../../stores";

    export let avatarPath = "";
    if (avatarPath === "") {
        avatarPath = "./avatars/default.png"
    }
    export let firstName = "";
    export let lastName = "";
    export let userID = "";

    async function addChatToBottom(targetID) {
        console.log("Target ID:", targetID)
        const chatContainer = document.getElementById('bottomChatContainer')
        if (!chatContainer) {
            console.error("Couldn't getElementById: #bottomChatContainer")
            return
        }
        // IF CHECK IF CHAT IS ALREADY THERE IF SO, return nil

        // Check if there is a chat ID between current WS/Client & targetUserID if not then request to create one 
        // return the chat ID 
        try {
            
            const chatID = await sendDataRequest({type: "getChatID", data:"I want some chatID data please!", id: $userInfo.id, targetid: targetID})
            console.log("i got the chatID:", chatID)
        
        } catch (error) {
            console.error("Error receiving chat ID:", error);
        }

        // Create the chat module


        console.log("hello!", firstName)
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="user" on:click={addChatToBottom(userID)}>
    <div class="profilePictureWrapper">
        <img src={avatarPath} alt={userID}>
    </div>

    <div class="usernameWrapper">
        <h2 class="username" style="margin: 0;">{firstName} {lastName}</h2>
    </div>

    <div class="messageNotification">
        <!-- Slide/Appear svelte animation **TODO -->
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
        border: 2px solid #5f9313bd;
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