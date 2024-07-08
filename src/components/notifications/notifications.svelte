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

    function clearNotifications() {
        sendMessage(
            JSON.stringify({ type: "clearNotif", fromid: $userInfo.id }),
        );
        notificationList = [];
        notifications.set(notificationList);
    }

    function deleteNotif(notifID) {}

    function handleNotificationClick(notificationID) {
        console.log("OU :D", notificationID);
        fetch("/markAsSeen", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: notificationID,
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

    async function updateFollowRequest(action, target) {
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
                type: "acceptedFollow",
                targetid: target, // see kes requesti saatis
                fromid: $userInfo.id, // see kes accept vajutas
                data: String,
            };

            console.log(messageData, "tererererer");
            messageData.data = "acceptedFollow_" + $userInfo.id.toString();
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
            {#each notificationList as notification}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li
                    id={notification.id}
                    on:click|once={() =>
                        handleNotificationClick(notification.id)}
                    class:clicked={notification.seen}
                >
                    {#if notification.data === undefined}
                        {notification.content}
                    {:else}
                        {notification.data}
                    {/if}
                    {#if notification.type === "followRequest"}
                        <button
                            on:click={() =>
                                updateFollowRequest(1, notification.fromID)}
                            >Accept</button
                        >
                        <button
                            on:click={() =>
                                updateFollowRequest(-1, notification.fromID)}
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
</style>
