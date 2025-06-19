let socket: WebSocket | null = null;

export function connectWebSocket() {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        socket = new WebSocket('ws://localhost:8080/ws');

        socket.onopen = () => console.log('Connected to WebSocket server');
        socket.onmessage = (event) => console.log('Message from server:', event.data);
        socket.onerror = (error) => console.error('WebSocket Error:', error);
        socket.onclose = () => console.log('WebSocket connection closed');
    }
}

export function sendWebSocketMessage(type: string, data: object = {}) {
    if (!localStorage.getItem('playerID')) {
        localStorage.setItem('playerID', crypto.randomUUID());
    }
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type, ...data, playerID: localStorage.getItem('playerID') }));
    } else {
        console.warn('WebSocket not connected');
    }
}
