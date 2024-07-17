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
        tabMap,
        isTypingStore,
        userInfo,
        allUsers,
        IMAGE_URL,
        allowedTabAmount,
    } from "../../stores";
    import { removeFromActiveChat } from "../../utils";
    import Message from "./message.svelte";
    import MinimizeChat from "../icons/minimizeChat.svelte";
    import CloseChat from "../icons/closeChat.svelte";
    import Chatbox from "./chatbox.svelte";

    $: users = $allUsers;
    //const tabMap = new Map();
    let firstTwoTabs = [];
    $: specialTabs = [];    
    let specialTabsOpen = false;

    $: if ($chatTabs.length >= 0) {
        firstTwoTabs = $chatTabs.slice(-$allowedTabAmount);
        specialTabs = $chatTabs.slice(0, -$allowedTabAmount);
        console.log("chatTabs:", $chatTabs);
        console.log("firstTwo:", firstTwoTabs);
        console.log("specialtabs:", specialTabs);

        $tabMap.forEach((unused, userID) => {
            if (!$chatTabs.some((tab) => tab.userID === userID)) {
                $tabMap.delete(userID);
            }
        });
        console.log("$tabMap", $tabMap);

        firstTwoTabs.forEach((tab) => {
            if (!$tabMap.has(tab.userID)) {
                addChatToBottom(
                    tab.userID,
                    tab.firstName,
                    tab.lastName,
                    tab.avatarPath,
                    tab.isGroup,
                    tab.groupChatID,
                );
                $tabMap.set(tab.userID, tab.firstName + " " + tab.lastName);
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
        // console.log("---GROUPS---");
        // console.log(targetID);
        // console.log(isGroup);
        // Incase of Groups the datatype is string with prefix 040
        if (isGroup === true) {
            let realid = parseInt(targetID.slice(12)); //sest safari on taun:D.
            const chatBox = new Chatbox({
                target: chatContainer,
                props: {
                    isGroup: true,
                    isFirstLoad: true,
                    userID: targetID,
                    chatID: realid,
                    userName: firstName + " " + lastName,
                    AvatarPath: "",
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
            return currentTabs.slice(-2);
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
        if (clickedChatIndex === -1) return;

        $chatTabs.push(...$chatTabs.splice(clickedChatIndex, 1));
        console.log($chatTabs);
        const removedUserID = $chatTabs[$chatTabs.length - 3].userID;
        removeFromActiveChat(event, "openChat", removedUserID);

        /* const clickedChatIndex = specialTabs.findIndex(
            (tab) => tab.userID === clickedUserID,
        );
        // Chat not found in specialTabs
        if (clickedChatIndex === -1) return;

        const [clickedChat] = specialTabs.splice(clickedChatIndex, 1);
        console.log(clickedChat);
        const lastFirstTwoTab = firstTwoTabs.shift();
        specialTabs.unshift(lastFirstTwoTab);

        firstTwoTabs.push(clickedChat);

        chatTabs.set([...firstTwoTabs, ...specialTabs]);
        const userID = lastFirstTwoTab.userID;
        removeFromActiveChat(event, "openChat", userID);
        $tabMap.delete(userID);*/
    }
</script>

<div class="specialTabContainer">
    {#if specialTabs.length > 0}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
            class="special-tab-preview"
            class:hidden={specialTabsOpen}
            on:click={() => (specialTabsOpen = true)}
        >
            <button class="red-button">{specialTabs.length}</button>
        </div>
        {#if specialTabsOpen}
            <div class="special-tab">
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div class="header">
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
                </div>
                {#each specialTabs as tab}
                    <div class="user" on:click={openChat(tab.userID)}>
                        <div class="profilePictureWrapper">
                            <img
                                src="{IMAGE_URL}{tab.avatarPath}"
                                alt="avatar"
                            />
                        </div>
                        <div class="usernameWrapper">
                            <h2 class="username" style="margin: 0;">
                                {tab.firstName}
                                {tab.lastName}
                            </h2>
                        </div>
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <div
                            class="close-chat"
                            on:click={deleteSingleChat(tab.userID)}
                        >
                            <CloseChat />
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    {/if}
</div>

<style>
    .specialTabContainer {
        transform: translatey(-180px);
    }
    .red-button {
        transform: translatey(+180px);

        cursor: pointer;
        background-color: red;
        color: white;
        border: none;
        border-radius: 50%;
        width: 25px;
        height: 25px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    :root {
        --chatWidth: 264px;
        --chatPreviewH: 40px;
        --chatFullH: 400px;
    }
    .special-tab {
        margin-right: 6px;
        margin-left: 6px;
        flex-direction: column;
        height: var(--chatFullH);
        width: var(--chatWidth);
        background-color: black;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.125);
        overflow-y: scroll;
        scrollbar-width: thin;
        scrollbar-color: greenyellow #011;
        overflow-x: none;
    }

    .minimize-tab,
    .close-chat {
        cursor: pointer;
        margin: 5px;
        align-items: right;
    }

    .header {
        display: flex;
        justify-content: space-between;
        width: var(--chatWidth);
        height: 60px;
    }

    .user {
        user-select: none;
        cursor: pointer;
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 3%;
        height: 42px;
        border-radius: 5px;
        border: 1px solid rgb(145, 145, 145);
        background-color: rgba(17, 25, 40, 0.75);
    }

    .user img {
        margin-right: 10px;
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
        background-color: rgba(17, 25, 40, 0.75);
    }
    .username {
        font-size: medium;
    }

    .hidden {
        display: none;
    }
</style>
