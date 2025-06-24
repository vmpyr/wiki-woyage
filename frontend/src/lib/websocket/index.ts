import { goto } from '$app/navigation';

let socket: WebSocket | null = null;

type Message = {
	type: string;
	[key: string]: any;
};

type MessageHandler = (msg: Message) => void;

const clientSideHandlers: Record<string, MessageHandler> = {
	"lobby_joined": (msg) => {
		if (msg.payload.lobbyID) {
            localStorage.setItem('username', msg.payload.username);
			localStorage.setItem('lobbyID', msg.payload.lobbyID);
			goto(`/lobby/${msg.payload.lobbyID}`);
		}
	},
};

export function connectWebSocket(username: string, lobbyID: string) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        socket = new WebSocket(`ws://localhost:8080/ws?username=${encodeURIComponent(username)}&lobbyID=${encodeURIComponent(lobbyID)}`);

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
                console.error('Invalid message from server:', event.data);
            }
        };

        socket.onerror = (error) => console.error('WebSocket Error:', error);

        socket.onclose = () => console.log('WebSocket connection closed');
    }
}
