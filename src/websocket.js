import { writable } from 'svelte/store';

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
        console.log("Recieved message:", response)

        if (pendingRequests[response.type]) {
            const { resolve, timeout } = pendingRequests[response.type];
            clearTimeout(timeout);
            resolve(response.data);
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

// Async data request through websocket
export const sendDataRequest = (message) => {
    return new Promise((resolve, reject) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            console.log("Sending WS data request:", message);
            socket.send(JSON.stringify(message));
            // Timeout for request
            const timeout = setTimeout(() => {
                reject(new Error(`WS sendDataRequest timeout (waiting for response ${message.type})`));
                delete pendingRequests[message.type];
            }, 5000); // 5 sec
            // Store resolve and timeout for this request
            pendingRequests[message.type] = { resolve, timeout };
        } else {
            reject(new Error('Websocket is not open'));
        }
    });
};