import { writable, get } from 'svelte/store';
import { InsertNewMessage } from './utils';
import { onlineUserStore, lastMsgStore, allUsers, chatNotifStore } from './stores';

export const messages = writable([]);
let socket;

// Map to store pending requests
const pendingRequests = {};

export const connect = (username) => {
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('WebSocket is connected');
        sendMessage(JSON.stringify({ type: "login", data: "", username:username }));
    };

    socket.onmessage = (event) => {
        const response = JSON.parse(event.data);
        // console.log("Recieved message:", response)

        // Update allUsers store
        if (response.type === "allUsers") {
            allUsers.set(response.allUsers)

        }

        // Update unseenMsgStore
        if (response.type === "chatNotifStore") {
            chatNotifStore.set(response.chatNotif)
        }

        // Update lastMsgs for userID on store
        if (response.type === "lastMsgStore") {
            lastMsgStore.set(response.lastMsgStore)
        }
        // Handle new incoming Message
        if (response.type === "newMessage") {
            InsertNewMessage(response)
        }
        // Update online users on store
        if (response.type === "onlineUsers") {
            onlineUserStore.set(response.onlineUsers)
        }

        if (pendingRequests[response.type]) {
            const { resolve, timeout } = pendingRequests[response.type];
            clearTimeout(timeout);
            resolve(response);
            // Remove it from pending req
            delete pendingRequests[response.type];
        } else {
            // Update messages store if its not a response to a request
            messages.update(msgs => [...msgs, event.data]);
        }
    };

    socket.onclose = () => {
        console.log('WebSocket is closed');
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
};

export const sendMessage = (message) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("Sending message:", message)
        socket.send(message);
    }
};

export const sendDataRequest = (request) => {
    return new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
            delete pendingRequests[request.type];
            reject(new Error('Request timed out'));
        }, 5000); // Timeout after 5 seconds

        pendingRequests[request.type] = { resolve, timeout };
        sendMessage(JSON.stringify(request));
    });
};