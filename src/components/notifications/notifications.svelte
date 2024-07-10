<script>
    import { notifications } from "../../websocket.js";
    import { onMount } from "svelte";
    import { userInfo } from "../../stores";
    import { sendMessage } from "../../websocket.js";

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

    $: sortedNotifications = notificationList.slice().sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));


    function clearNotifications() {
        sendMessage(
            JSON.stringify({ type: "clearNotif", fromid: $userInfo.id }),
        );
        notificationList = [];
        notifications.set(notificationList);
    }

    function clearSingleNotification(notifID, targetid) {
        sendMessage(
            JSON.stringify({ type: "clearSingleNotif", data: notifID.toString(), targetid: targetid}),
        );
        console.log("clearing single notification", notifID, targetid);
        removeNotification(notifID);
    }

    function handleNotificationClick(notificationID) {
        console.log("OU :D", notificationID);
        fetch("/markAsSeen", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ notificationID }),
        }).catch((err) => console.error(err));

        notificationList = notificationList.map((notification) => {
            if (notification.id === notificationID) {
                return {
                    ...notification,
                    seen: true,
                };
            }
            return notification;
        });
    }

    function removeNotification(notifID) {
        notificationList = (notificationList.filter(notification => notification.id !== notifID));
        notifications.set(notificationList);
}

    async function updateFollowRequest(action, target, notifID) {
        console.log("updateFollowRequest:", action, target);
        try {
            const response = await fetch("/api/followers", {
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
                messageData.type = "cancelRequest"
                messageData.data = $userInfo.id.toString();
            } else if (action === 1) {
                messageData.type = "acceptedFollow"
                messageData.data = "acceptedFollow_" + $userInfo.id.toString();
            }

            console.log(messageData, "tererererer");
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
            {#each sortedNotifications  as notification}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li
                    id={notification.id}
                    on:click|once={() =>
                        handleNotificationClick(notification.id)}
                    class:clicked={notification.seen}
                >
                    <button class="close-button" on:click={() => clearSingleNotification(notification.id, notification.fromID)}>X</button>

                    {#if notification.content !== undefined}
                    {notification.content}
                    {/if}

                    {#if notification.type === "followRequest"}
                        <button
                            on:click={() => { updateFollowRequest(1, notification.fromID, notification.id); removeNotification(notification.id);}}>
                            Accept
                            </button>
                        <button
                            on:click={() =>
                                { updateFollowRequest(-1, notification.fromID, notification.id); removeNotification(notification.id);}}
                            >Decline</button
                        >
                    {/if}
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
    }

    button {
        margin-left: 10px;
        padding: 5px 10px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        background-color: green;
    }

    button:hover {
        background-color: #ddd;
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
        position: absolute;
        top: 5px;
        right: 5px;
        background: transparent;
        border: none;
        font-size: 1.2em;
        cursor: pointer;
    }
</style>
