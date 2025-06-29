import { GOLANG_WS_URL, MAX_RECONNECT_ATTEMPTS, RECONNECT_DELAY } from '$lib';
import type { Response } from './responsehandler';
import { responseHandlers } from './responsehandler';

export let socket: WebSocket | null = null;

let isIntentionalDisconnect = false;
let reconnectAttempts = 0;
let reconnectTimeout: number | null = null;

export function connectWebSocket(username: string, lobbyID: string) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        let clientID = localStorage.getItem('clientID');
        if (!clientID) {
            clientID = crypto.randomUUID();
            localStorage.setItem('clientID', clientID);
        }
        const params = new URLSearchParams({ username, lobbyID, clientID });
        socket = new WebSocket(`${GOLANG_WS_URL}/ws?${params.toString()}`);

        socket.onopen = () => {
            console.log('Connected to WebSocket server');
            reconnectAttempts = 0;
            setupGracefulTermination();
        };

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

        socket.onclose = (event) => {
            console.log('WebSocket connection closed', event.code, event.reason);

            window.removeEventListener('beforeunload', handlePageUnload);
            window.removeEventListener('pagehide', handlePageUnload);

            if (!isIntentionalDisconnect && event.code !== 1000) {
                attemptReconnect(username, lobbyID);
            }
        };
    }
}

function setupGracefulTermination() {
    window.addEventListener('beforeunload', handlePageUnload);
    window.addEventListener('pagehide', handlePageUnload);
}

export function sendDisconnectEvent(intentional: boolean) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        try {
            socket.send(
                JSON.stringify({
                    type: 'disconnect',
                    handler: 'orchestrator',
                    payload: {}
                })
            );
            isIntentionalDisconnect = intentional;
        } catch (error) {
            console.error('Failed to send disconnect event:', error);
        }
    }

    if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        reconnectTimeout = null;
    }
    if (socket) {
        socket.close();
        socket = null;
    }
}

function attemptReconnect(username: string, lobbyID: string) {
    if (isIntentionalDisconnect || reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
        console.log('Max reconnection attempts reached or intentional disconnect');
        return;
    }

    reconnectAttempts++;
    console.log(`Attempting to reconnect... (${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})`);

    reconnectTimeout = setTimeout(() => {
        connectWebSocket(username, lobbyID);
    }, RECONNECT_DELAY * reconnectAttempts);
}

function handlePageUnload() {
    sendDisconnectEvent(false);
}
