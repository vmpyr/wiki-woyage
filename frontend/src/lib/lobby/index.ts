import { GOLANG_HTTP_URL } from '$lib';
import { socket } from '$lib/websocket';

type ClientInfo = {
    username: string;
    lobbyID: string;
};

async function fetchClientInfo(): Promise<ClientInfo | null> {
    const clientID = localStorage.getItem('clientID');
    if (!clientID) return null;
    try {
        const res = await fetch(
            `${GOLANG_HTTP_URL}/api/clientinfo?clientID=${encodeURIComponent(clientID)}`
        );
        if (!res.ok) return null;
        const data = await res.json();
        return {
            username: data.username,
            lobbyID: data.lobbyID
        };
    } catch (err) {
        console.error('Failed to fetch client info:', err);
        return null;
    }
}

export async function getClientInfo(): Promise<ClientInfo | null> {
    return await fetchClientInfo();
}

// TODO: I am just testing, generally player will not explicitly request admin status (imagine though)
export function requestAdminStatus() {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'am_i_admin', handler: 'lobby', payload: {} }));
    }
}
