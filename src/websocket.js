import { writable } from 'svelte/store';

export const messages = writable([]);
let socket;

export const connect = (username) => {
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('WebSocket is connected');
        sendMessage(JSON.stringify({ type: "login", data: "", username:username }));
    };

    socket.onmessage = (event) => {
        messages.update(msgs => [...msgs, event.data]);
    };

    socket.onclose = () => {
        console.log('WebSocket is closed');
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
};

export const sendMessage = (message) => { // message format { type: "type", data: "data", username:username }
    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("Sending message:", message)
        socket.send(message);
    }
};
