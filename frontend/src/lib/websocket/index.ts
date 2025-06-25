import { goto } from '$app/navigation';
import { GOLANG_WS_URL } from '$lib';

let socket: WebSocket | null = null;

type Message = {
    type: string;
    payload: Record<string, unknown>;
};

type MessageHandler = (msg: Message) => void;

const clientSideHandlers: Record<string, MessageHandler> = {
    lobby_joined: (msg) => {
        const lobbyID = msg.payload.lobbyID;
        if (typeof lobbyID === 'string') {
            goto(`/lobby/${lobbyID}`);
        } else {
            console.warn('lobbyID is not a string:', lobbyID);
        }
    }
};

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
                const msg: Message = JSON.parse(event.data);
                const clientSideHandler = clientSideHandlers[msg.type];
                if (clientSideHandler) {
                    clientSideHandler(msg);
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
