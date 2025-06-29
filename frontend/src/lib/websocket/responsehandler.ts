import { goto } from '$app/navigation';
import { playerList } from '$lib/stores';

export type Response = {
    type: string;
    response: Record<string, unknown>;
};

type ResponseHandler = (msg: Response) => void;

export const responseHandlers: Record<string, ResponseHandler> = {
    lobby_joined: (msg) => {
        const lobbyID = msg.response.lobbyID;
        if (typeof lobbyID === 'string') {
            goto(`/lobby/${lobbyID}`);
        } else {
            console.warn('lobbyID is not a string:', lobbyID);
        }
    },

    player_list: (msg) => {
        const playerListResponse = msg.response.players;
        if (Array.isArray(playerListResponse)) {
            const playerNamesWithStatus = playerListResponse.map((player: string) => [
                player,
                'joined'
            ]);
            playerList.set(playerNamesWithStatus);
        }
    },

    new_player: (msg) => {
        const newPlayer = msg.response.username;
        if (typeof newPlayer === 'string') {
            playerList.update((list) => [...list, [newPlayer, 'joined']]);
        }
    },

    player_left: (msg) => {
        const leavingPlayer = msg.response.username;
        if (typeof leavingPlayer === 'string') {
            playerList.update((list) => [...list, [leavingPlayer, 'left']]);
        }
    },

    // just testing
    am_i_admin: (msg) => {
        const isAdmin = msg.response.isAdmin;
        if (typeof isAdmin !== 'boolean') {
            console.warn('isAdmin is not a boolean');
            return;
        }
        alert(`You are ${isAdmin ? 'the admin' : 'not the admin'}`);
    }
};
