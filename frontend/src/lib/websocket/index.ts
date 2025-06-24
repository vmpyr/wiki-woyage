let socket: WebSocket | null = null;

export function connectWebSocket(username: string, lobbyID: string) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        socket = new WebSocket(`ws://localhost:8080/ws?username=${encodeURIComponent(username)}&lobbyID=${encodeURIComponent(lobbyID)}`);

        socket.onopen = () => console.log('Connected to WebSocket server');
        socket.onmessage = (event) => console.log('Message from server:', event.data);
        socket.onerror = (error) => console.error('WebSocket Error:', error);
        socket.onclose = () => console.log('WebSocket connection closed');
    }
}
