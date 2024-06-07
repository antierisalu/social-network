import { writable } from 'svelte/store';

export const messages = writable([]);
let socket;

export const connect = () => {
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('WebSocket is connected');
    };

    socket.onmessage = (event) => {
        console.log(event.data)
        messages.update(msgs => [...msgs, event.data]);
        console.log(messages)
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
