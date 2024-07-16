import { writable, get } from "svelte/store";
import { InsertNewMessage, bellNotif } from "./utils";
import { onlineUserStore, lastMsgStore, allUsers, chatNotifStore, groupChatNotifStore, setTyping, removeTyping } from './stores';

export const messages = writable([]);
export const notifications = writable([]);

let socket;
const audio = new Audio("notification.mp3");
audio.volume = 0.1;
let originalTitle = document.title;
let titleTimeout;

function playSound(){
    audio.pause()
    audio.currentTime = 0
    audio.play();
}

// Map to store pending requests
const pendingRequests = {};

export const connect = (username) => {
    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
        sendMessage(
            JSON.stringify({ type: "login", data: "", username: username })
        );
    };

    if ("Notification" in window) {
        Notification.requestPermission().then(function (permission) {
            if (permission === "granted") {
                console.log("Notifications allowed");
            } else {
                console.log("Notifications denied");
            }
        });
    } else {
        alert("This browser does not support desktop notification");
    }

    socket.onmessage = (event) => {

        const response = JSON.parse(event.data);
        // console.log("Recieved message:", response)

        switch (response.type) {
            case "newMessage":
                InsertNewMessage(response);
                removeTyping(response.fromUserID)
                break;
            case "newGroupMessage":
                InsertNewMessage(response, true);
                break;
            case "followRequest":
                updateTabTitle("New notification");
                console.log("YOU RECEIVED A NOTIFICATION");
                notifications.update((n) => [...n, response]);
                playSound();
                bellNotif();
                break;
            case "follow":
                updateTabTitle("New notification");
                console.log("YOU RECEIVED A NOTIFICATION");
                notifications.update((n) => [...n, response]);
                playSound();
                bellNotif();
                break;
            case "acceptedFollow":
                updateTabTitle("New notification");
                console.log("YOU RECEIVED A NOTIFICATION");
                notifications.update((n) => [...n, response]);
                playSound();
                bellNotif();
                break;
            case "isTyping" :
                setTyping(response.fromID)
                break;
            // Update allUsers store
            case "allUsers":
                allUsers.set(response.allUsers)
                break;
            // Update unseenMsgStore (PM)
            case "chatNotifStore":
                chatNotifStore.set(response.chatNotif)
                break;
            // Update unseenMsgStore (GM)
            case "groupChatNotifStore":
                console.log("groupChatNotifStore: ", response, response.chatNotif)
                if (response.chatNotif !== null) {
                    groupChatNotifStore.set(response.chatNotif)
                }
                break;
            // Update lastMsgs for userID on store
            case "lastMsgStore":
                lastMsgStore.set(response.lastMsgStore)
                break;
            case "onlineUsers" :
                onlineUserStore.set(response.onlineUsers)
                break;
            case "cancelRequest":
                notifications.update((n) => n.filter(notification => notification.id !== response.id));
            case "groupRequest":
                console.log("Group request recieved", response)
                notifications.update((n) => [...n, response]);
                break;
            case "groupInvite":
                notifications.update((n) => [...n, response]);
                console.log("Group request recieved", response)
                break;
        }

        if (pendingRequests[response.type]) {
            const { resolve, timeout } = pendingRequests[response.type];
            clearTimeout(timeout);
            resolve(response);
            // Remove it from pending req
            delete pendingRequests[response.type];
        } else {
            // Update messages store if it's not a response to a request
            messages.update((msgs) => [...msgs, event.data]);
        }
    };

    socket.onclose = () => {
        console.log("WebSocket is closed");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};

export const sendMessage = (message) => {
    // message format { type: "type", data: "data", username:username }
    if (socket && socket.readyState === WebSocket.OPEN) {
        // console.log("Sending message:", message);
        socket.send(message);
    }
};

export const sendDataRequest = (request) => {
    return new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
            delete pendingRequests[request.type];
            reject(new Error("Request timed out"));
        }, 5000);

        pendingRequests[request.type] = { resolve, timeout };
        sendMessage(JSON.stringify(request));
    });
};

function updateTabTitle(notification) {
    originalTitle = document.title;
    document.title = notification;

    function onVisibilityChange() {
        if (!document.hidden) {
            document.title = originalTitle;
            clearTimeout(titleTimeout);
            document.removeEventListener(
                "visibilitychange",
                onVisibilityChange
            );
        }
    }

    titleTimeout = setTimeout(() => {
        document.title = originalTitle;
        document.removeEventListener("visibilitychange", onVisibilityChange);
    }, 5000);

    document.addEventListener("visibilitychange", onVisibilityChange);
}
