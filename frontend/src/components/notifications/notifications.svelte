<script>
    import { notifications } from "../../websocket.js";
    import { onMount } from "svelte";
    import { userInfo, activeTab, API_URL } from "../../stores";
    import { sendMessage } from "../../websocket.js";
    import { selectUser } from "../../utils.js";

    let notificationList = [];

    onMount(() => {
        const unsubscribe = notifications.subscribe((newNotifications) => {
            console.log("newnotifications: ", newNotifications);
            notificationList = newNotifications;
        });

        return () => {
            unsubscribe();
        };
    });

    $: sortedNotifications = notificationList
        .slice()
        .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));

    function clearNotifications() {
        sendMessage(
            JSON.stringify({ type: "clearNotif", fromid: $userInfo.id }),
        );
        notificationList = [];
        notifications.set(notificationList);
    }

    function clearSingleNotification(notifID, targetid) {
        sendMessage(
            JSON.stringify({
                type: "clearSingleNotif",
                data: notifID.toString(),
                targetid: targetid,
            }),
        );
        console.log("clearing single notification", notifID, targetid);
        removeNotification(notifID);
    }

    function handleNotificationClick(notification) {
        console.log("OU :D", notification.id);
        console.log("notification: ", notification);
        let notificationID = notification.id;
        console.log("notifID :D", notificationID);
        fetch(`${API_URL}/markAsSeen`, {
            credentials: "include",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ notificationID }),
        }).catch((err) => console.error(err));

        notificationList = notificationList.map((notification) => {
            if (notification.id === notification.id) {
                return {
                    ...notification,
                    seen: true,
                };
            }
            return notification;
        });
        activeTab.set("Profile");
        selectUser(notification.fromID);
    }

    function removeNotification(notifID) {
        notificationList = notificationList.filter(
            (notification) => notification.id !== notifID,
        );
        notifications.set(notificationList);
    }

    async function updateFollowRequest(action, target, notifID) {
        console.log("updateFollowRequest:", action, target);
        try {
            const response = await fetch(`${API_URL}/followers`, {
                credentials: "include",
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ action: action, target: target }),
            });

            let userData = await response.json(); //returns who initiated follow change
            var messageData = {
                type: String,
                targetid: target, // see kes requesti saatis
                fromid: $userInfo.id, // see kes accept vajutas
                data: String,
                notificationid: notifID,
            };

            if (action === -1) {
                messageData.type = "declinedFollow";
                messageData.data = "declinedFollow_" + $userInfo.id.toString();
            } else if (action === 1) {
                messageData.type = "acceptedFollow";
                messageData.data = "acceptedFollow_" + $userInfo.id.toString();
            }

            sendMessage(JSON.stringify(messageData));

            console.log(userData);
        } catch (error) {
            console.error(
                "Error sending update follow request: ",
                error.message,
            );
        }
    }
</script>

<main>
    <h1>Notifications ({notificationList.length})</h1>
    {#if notificationList.length > 0}
        <ul>
            {#each sortedNotifications as notification}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li id={notification.id} class:clicked={notification.seen}>
                    <div class="close-btn">
                        <button
                            class="close-button"
                            on:click={() =>
                                clearSingleNotification(
                                    notification.id,
                                    notification.fromID,
                                )}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                x="0px"
                                y="0px"
                                width="100"
                                height="100"
                                viewBox="0 0 50 50"
                            >
                                <path
                                    d="M 7.71875 6.28125 L 6.28125 7.71875 L 23.5625 25 L 6.28125 42.28125 L 7.71875 43.71875 L 25 26.4375 L 42.28125 43.71875 L 43.71875 42.28125 L 26.4375 25 L 43.71875 7.71875 L 42.28125 6.28125 L 25 23.5625 Z"
                                ></path>
                            </svg>
                        </button>
                    </div>
                    <div
                        class="notification-content"
                        on:click|once={() =>
                            handleNotificationClick(notification)}
                    >
                        {#if notification.content !== undefined}
                            {notification.content}
                        {/if}
                    </div>

                    <div class="action-buttons">
                        {#if notification.type === "followRequest"}
                            <button
                                on:click={() => {
                                    updateFollowRequest(
                                        1,
                                        notification.fromID,
                                        notification.id,
                                    );
                                    removeNotification(notification.id);
                                }}
                            >
                                Accept
                            </button>
                            <button
                                on:click={() => {
                                    updateFollowRequest(
                                        -1,
                                        notification.fromID,
                                        notification.id,
                                    );
                                    removeNotification(notification.id);
                                }}
                            >
                                Decline
                            </button>
                        {/if}
                    </div>
                </li>
            {/each}
        </ul>
        <button on:click={clearNotifications}>Clear Notifications</button>
    {/if}
    {#if notificationList.length === 0}
        <p>No notifications</p>
    {/if}
</main>

<style>
    main {
        padding: 20px;
        font-family: Arial, sans-serif;
    }

    h1 {
        font-size: 24px;
        margin-bottom: 10px;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li {
        background-color: #f4f4f4;
        padding: 10px;
        margin-bottom: 5px;
        border-radius: 4px;
        position: relative;
        display: flex;
        flex-direction: column;
    }

    .notification-content {
        margin-bottom: 10px;
        padding-left: 24px;
        padding-right: 24px;
    }

    button {
        border: none;
        border-radius: 4px;
        cursor: pointer;
        background-color: green;
        color: white;
    }

    button:hover {
        background-color: #ddd;
        color: black;
    }

    p {
        color: #888;
    }

    li.clicked {
        background-color: lightgray;
        -webkit-user-select: none;
        -moz-user-select: none;
        user-select: none;
    }

    li:hover {
        cursor: pointer;
    }

    .close-button {
        font-size: 12px;
        background-color: transparent;
        color: black;
        margin: 0px;
    }

    .close-button svg {
        width: 18px;
        height: 18px;
    }

    .close-btn {
        margin-bottom: -16px;
        display: flex;
        justify-content: flex-end;
        height: auto;
    }

    .close-btn::hover {
        color: black;
    }

    .action-buttons {
        display: flex;
        justify-content: center;
        gap: 10px;
    }

    .action-buttons button {
        margin: 0;
    }
</style>
