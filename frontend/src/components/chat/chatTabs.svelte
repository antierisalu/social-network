<script>
    import { writable } from "svelte/store";
    import MsgNotification from "../icons/msgNotification.svelte";
    import {
        connect,
        sendMessage,
        messages,
        sendDataRequest,
    } from "../../websocket";
    import { get } from "svelte/store";
    import {
        activeTab,
        chatTabs,
        isTypingStore,
        userInfo,
        allUsers,
        IMAGE_URL
    } from "../../stores";
    import { removeFromActiveChat } from "../../utils";
    import Message from "./message.svelte";
    import MinimizeChat from "../icons/minimizeChat.svelte";
    import CloseChat from "../icons/closeChat.svelte";
    import Chatbox from "./chatbox.svelte";

    $: users = $allUsers;
    const tabMap = new Map();
    let firstTwoTabs = [];
    $: specialTabs = [];
    let specialTabsOpen = false;

    $: if ($chatTabs.length > 0) {
        firstTwoTabs = $chatTabs.slice(0, 2);
        specialTabs = $chatTabs.slice(2);
        console.log("chatTabs:", $chatTabs);
        console.log("firstTwo:", firstTwoTabs);
        console.log("specialtabs:", specialTabs);

        tabMap.forEach((unused, userID) => {
            if (!$chatTabs.some((tab) => tab.userID === userID)) {
                tabMap.delete(userID);
            }
        });
        console.log("tabMap", tabMap);

        firstTwoTabs.forEach((tab) => {
            if (!tabMap.has(tab.userID)) {
                addChatToBottom(
                    tab.userID,
                    tab.firstName,
                    tab.lastName,
                    tab.avatarPath,
                    tab.isGroup,
                    tab.groupChatID,
                );
                tabMap.set(tab.userID, tab.firstName + " " + tab.lastName);
            }
        });
    }

    async function addChatToBottom(
        targetID,
        firstName,
        lastName,
        avatarPath,
        isGroup,
        amogus,
    ) {
        if (targetID === $userInfo.id) {
            console.log("cant message yourself!");
            return;
        }

        const chatContainer = document.getElementById("bottomChatContainer");
        if (!chatContainer) {
            console.error("Couldn't getElementById: #bottomChatContainer");
            return;
        }

        // GROUPS
        console.log("---GROUPS---");
        console.log(targetID);
        console.log(isGroup);
        // Incase of Groups the datatype is string with prefix 040
        if (isGroup === true) {
            let realid = parseInt(targetID.slice(12)) //sest safari on taun:D.
            const chatBox = new Chatbox({
                target: chatContainer,
                props: {
                    isGroup: true,
                    isFirstLoad: true,
                    userID: targetID,
                    chatID: realid,
                    userName: firstName + " " + lastName,
                    AvatarPath: avatarPath,
                },
            });

            return;
        } else {
            console.log("this is not a grouP!");
        }
        // GROUPS

        // Check if there is a chat ID between current WS/Client & targetUserID if not then request to create one
        // return the chat ID
        try {
            console.log("see on priv chattieledledle");
            const response = await sendDataRequest({
                type: "getChatID",
                data: "",
                id: $userInfo.id,
                targetid: targetID,
            });
            var amogus = response.chatID;
            const targetUserData = users.find((user) => user.ID === targetID);
            if (!targetUserData) {
                console.log(
                    "Failed to get target user's data from store/allUsers",
                );
            }

            // To not open more than one chat tabs with same user

            const existingChatBox = chatContainer.querySelector(
                `div[chatid="${amogus}"]`,
            );
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
                    isGroup: false,
                    isFirstLoad: true,
                    userID: targetID,
                    chatID: amogus,
                    userName: firstName + " " + lastName,
                    AvatarPath: targetUserData.Avatar,
                },
            });
        } catch (error) {
            console.error("Error receiving chat ID:", error);
        }
        // const chatBody = chatContainer.querySelector(`div[chatid="${chatID}"]`)
    }

    function deleteAllChats() {
        chatTabs.update((currentTabs) => {
            return currentTabs.slice(0, 2);
        });
        specialTabsOpen = false;
    }

    function deleteSingleChat(userID) {
        chatTabs.update((currentTabs) => {
            return currentTabs.filter((id) => id.userID !== userID);
        });
        if (specialTabs.length == 1) {
            specialTabsOpen = false;
        }
    }

    function openChat(clickedUserID) {
        const clickedChatIndex = specialTabs.findIndex(
            (tab) => tab.userID === clickedUserID,
        );
        // Chat not found in specialTabs
        if (clickedChatIndex === -1) return;

        const [clickedChat] = specialTabs.splice(clickedChatIndex, 1);

        const lastFirstTwoTab = firstTwoTabs.shift();
        specialTabs.unshift(lastFirstTwoTab);

        firstTwoTabs.push(clickedChat);

        chatTabs.set([...firstTwoTabs, ...specialTabs]);

        const userID = lastFirstTwoTab.userID;
        removeFromActiveChat(event, "openChat", userID);
        tabMap.delete(userID);
    }
</script>

{#if specialTabs.length > 0}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="special-tab-preview" on:click={() => (specialTabsOpen = true)}>
        <p>chats opened: {specialTabs.length}</p>
    </div>
    {#if specialTabsOpen}
        <div class="special-tab">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div
                class="minimize-tab"
                on:click={() => (specialTabsOpen = false)}
            >
                <MinimizeChat />
            </div>
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="close-chat" on:click={deleteAllChats}>
                <CloseChat />
            </div>
            {#each specialTabs as tab}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div class="close-chat" on:click={deleteSingleChat(tab.userID)}>
                    <CloseChat />
                </div>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div class="user" on:click={openChat(tab.userID)}>
                    <img src={IMAGE_URL}{tab.avatarPath} alt="avatar" />
                    <p>{tab.firstName} {tab.lastName}</p>
                    <!-- <div class="avatar {(isOnline) ? 'online' : 'offline'}">
                            <img src={tab.avatarPath} alt={tab.userID} class="{(isOnline) ? '' : 'avatar-grayscale'}"> -->
                </div>
            {/each}
        </div>
    {/if}
{/if}

<style>
    img {
        max-width: 20px;
    }
</style>
