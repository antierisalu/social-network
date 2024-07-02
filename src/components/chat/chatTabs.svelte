<script>
    import { writable } from "svelte/store"; 
    import MsgNotification from "../icons/msgNotification.svelte";
    import { connect, sendMessage, messages, sendDataRequest } from "../../websocket";
    import { get } from "svelte/store";
    import { activeTab, chatTabs, isTypingStore, userInfo, allUsers } from "../../stores";
    import Message from './message.svelte';
    import Chatbox from "./chatbox.svelte";
    
    $: users = $allUsers;
    const tabMap = new Map ()

    $: if ($chatTabs.length > 0) {
        console.log('chatTabs:',$chatTabs)
        const uniqueUserIDs = new Set();
        const uniqueTabs = $chatTabs.filter(tab => {
            const isUnique = !uniqueUserIDs.has(tab.userID);
            uniqueUserIDs.add(tab.userID);
            return isUnique;
        });

        const firstTwoTabs = uniqueTabs.slice(0, 2);
        const specialTabs = uniqueTabs.slice(2);
        console.log('firstTwo:', firstTwoTabs);
        console.log('specialtabs:', specialTabs);

        firstTwoTabs.forEach(tab => {
            if (!tabMap.has(tab.userID)) {
                addChatToBottom(tab.userID, tab.firstName, tab.lastName, tab.avatarPath);
                tabMap.set(tab.userID, true);
            }
        });

        specialTabs.forEach(tab => {
            if (!tabMap.has(tab.userID)) {
                buildSpecialTab(tab.userID, tab.firstName, tab.lastName, tab.avatarPath);
                tabMap.set(tab.userID, true);
            }
        });
    }



    async function addChatToBottom(targetID, firstName, lastName, avatarPath) {
        
        if (targetID === $userInfo.id) {
            console.log("cant message yourself!")
            return
        }

        const chatContainer = document.getElementById('bottomChatContainer')
        if (!chatContainer) {
            console.error("Couldn't getElementById: #bottomChatContainer")
            return
        }
        
        // Check if there is a chat ID between current WS/Client & targetUserID if not then request to create one 
        // return the chat ID
        try {
            const response = await sendDataRequest({type: "getChatID", data:"", id: $userInfo.id, targetid: targetID})
            var chatID = response.chatID;
            const targetUserData = users.find((user) => user.ID === targetID)
            if (!targetUserData) {
                console.log("Failed to get target user's data from store/allUsers")
            }

            // To not open more than one chat tabs with same user
        
            const existingChatBox = chatContainer.querySelector(`div[chatid="${chatID}"]`);
            if (existingChatBox) {
                console.log("Chat with this user is already open");
                return;
            }

            // Check for relationship type & pass it into prop (for followers+chats ***)
            // Create the chatBox module

            //saada see välja objektina arraysse aga võta chatID välja sest seda vb ei eksisteeri veel.
            const chatBox = new Chatbox({
                target: chatContainer,
                props: {
                    isFirstLoad: true,
                    userID: targetID,
                    chatID: chatID,
                    userName: (firstName + " " + lastName),
                    AvatarPath: targetUserData.Avatar,

                }
            });

            

        } catch (error) {
            console.error("Error receiving chat ID:", error);
        }
        // const chatBody = chatContainer.querySelector(`div[chatid="${chatID}"]`)

    }

    export function buildSpecialTab(targetID,firstName, lastName, avatarPath) {

    }

</script>






