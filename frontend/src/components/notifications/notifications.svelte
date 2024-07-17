<script>
    import { notifications } from "../../websocket.js";
    import { onMount } from "svelte";
    import { userInfo, activeTab, API_URL, groupSelected } from "../../stores";
    import { sendMessage } from "../../websocket.js";
    import { selectUser, getGroups } from "../../utils.js";

    import Close from "./closeNotif.svelte";
    import ActionButtons from "./actionButtons.svelte";

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
        if (
            notification.type === "groupInvite" ||
            notification.type === "groupRequest" ||
            notification.type === "acceptedGroupRequest"
        ) {
            activeTab.set("Groups");
            $groupSelected = parseInt(notification.link.split("_")[2]);
            getGroups();
            return;
        }
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
            var messageData = {
                type: String,
                targetid: target, // see kes requesti saatis
                fromid: $userInfo.id, // see kes accept vajutas
                data: String,
                notificationid: notifID,
            };

            if (action === -1) {
                messageData.type = "declinedRequest";
                messageData.data = "declinedRequest_" + $userInfo.id.toString();
            } else if (action === 1) {
                messageData.type = "acceptedFollow";
                messageData.data = "acceptedFollow_" + $userInfo.id.toString();
            }

            sendMessage(JSON.stringify(messageData));
        } catch (error) {
            console.error(
                "Error sending update follow request: ",
                error.message,
            );
        }
    }

    async function updateGroupRequest(action, target, notif) {
        console.log("updateGroupRequest:", action, target, notif);
        let groupID = parseInt(notif.link.split("_")[2]);
        try {
            const response = await fetch(`${API_URL}/joinGroup`, {
                credentials: "include",
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    action: action,
                    targetID: target,
                    groupID: groupID,
                }),
            });
            var messageData = {
                type: String,
                targetid: target, // see kes requesti saatis
                fromid: $userInfo.id, // see kes accept vajutas
                data: String,
                notificationid: notif.id,
                groupID: groupID,
            };

            if (action === -1) {
                messageData.type = "declinedRequest";
                messageData.data = "declinedRequest_" + $userInfo.id.toString();
            } else if (action === 1) {
                messageData.type = "acceptedGroupRequest";
                messageData.data =
                    "acceptedGroupRequest_" + $userInfo.id.toString() + "_" + groupID;
            }

            sendMessage(JSON.stringify(messageData));
            //TODO: saada websocket tagasi userile et request on t2idetud, handlei notif deletion
        } catch (error) {
            console.error(
                "Error sending update group request: ",
                error.message,
            );
        }
    }
    async function updateGroupInvite(action, notif, target) {
        console.log("updateGroupInvite:", action, notif);
        let groupID = parseInt(notif.link.split("_")[2]);
        try {
            const response = await fetch(`${API_URL}/joinGroup`, {
                credentials: "include",
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    action: action,
                    targetID: $userInfo.id,
                    groupID: groupID,
                }),
            });

            var messageData = {
                type: String,
                targetid: target, // see kes requesti saatis
                fromid: $userInfo.id, // see kes accept vajutas
                data: String,
                notificationid: notif.id,
                groupID: groupID,
            };

            if (action === -1) {
                messageData.type = "declinedRequest";
                messageData.data = "declinedRequest_" + $userInfo.id.toString();
            } else if (action === 1) {
                messageData.type = "acceptedGroupInvite";
                messageData.data =
                    "acceptedGroupInvite_" + $userInfo.id.toString();
            }
            console.log(messageData);
            sendMessage(JSON.stringify(messageData));
            //TODO: saada websocket tagasi userile et invite on t2idetud, handlei notif deletion
        } catch (error) {
            console.error(
                "Error sending update group request: ",
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
                    <Close
                        on:click={() =>
                            clearSingleNotification(
                                notification.id,
                                notification.fromID,
                            )}
                    />

                    <div
                        class="notification-content"
                        on:click|once={() =>
                            handleNotificationClick(notification)}
                    >
                        {#if notification.content !== undefined}
                            {notification.content}
                        {/if}
                    </div>

                    <ActionButtons
                        {notification}
                        on:updateremove={(event) => {
                            removeNotification(notification.id);
                            switch (notification.type) {
                                case "followRequest":
                                    updateFollowRequest(
                                        event.detail.action,
                                        notification.fromID,
                                        notification.id,
                                    );
                                    break;
                                case "groupRequest":
                                    updateGroupRequest(
                                        event.detail.action,
                                        notification.fromID,
                                        notification,
                                    );
                                    break;
                                case "groupInvite":
                                    updateGroupInvite(
                                        event.detail.action,
                                        notification,
                                        notification.fromID,
                                    );
                                    break;
                            }
                        }}
                    />
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
</style>
