import { GOLANG_WS_URL } from '$lib';
import type { Response } from './responsehandler';
import { responseHandlers } from './responsehandler';

export let socket: WebSocket | null = null;

export function connectWebSocket(username: string, lobbyID: string) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        let clientID = localStorage.getItem('clientID');
        if (!clientID) {
            clientID = crypto.randomUUID();
            localStorage.setItem('clientID', clientID);
        }
        const params = new URLSearchParams({ username, lobbyID, clientID });
        socket = new WebSocket(`${GOLANG_WS_URL}/ws?${params.toString()}`);

        socket.onopen = () => console.log('Connected to WebSocket server');

        socket.onmessage = (event) => {
            try {
                const msg: Response = JSON.parse(event.data);
                const responseHandler = responseHandlers[msg.type];
                if (responseHandler) {
                    responseHandler(msg);
                } else {
                    console.warn('No handler for message type:', msg.type, msg);
                }
            } catch (err) {
                console.error('Invalid message from server:', event.data, err);
            }
        };

        socket.onerror = (error) => console.error('WebSocket Error:', error);

        socket.onclose = () => console.log('WebSocket connection closed');
    }
}
