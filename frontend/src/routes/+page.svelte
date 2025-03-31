<script>
    import { onMount } from 'svelte';
    import { connectWebSocket, sendWebSocketMessage } from '$lib/websocket';

    let playerName = '';
    let lobbyID = '';

    onMount(() => {
        connectWebSocket();
    });

    function createLobby() {
        sendWebSocketMessage('create_lobby', { playerName });
    }

    function joinLobby() {
        if (lobbyID.trim() === '') return;
        sendWebSocketMessage('join_lobby', { lobbyID, playerName });
    }
</script>

<main class="flex h-screen flex-col items-center justify-center bg-gray-800 text-white">
    <h1 class="mb-4 text-3xl font-bold">Welcome to WikiWoyage!</h1>

    <div class="mb-4">
        <input
            type="text"
            id="playerName"
            bind:value={playerName}
            class="mt-2 w-64 rounded border border-gray-300 p-2"
            placeholder="Your name"
        />
        <input
            type="text"
            id="lobbyID"
            bind:value={lobbyID}
            class="mt-2 w-64 rounded border border-gray-300 p-2"
            placeholder="Lobby ID (for joining only)"
        />
        <button
            class="mt-2 rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
            on:click={createLobby}
        >
            Create Lobby
        </button>
        <button
            class="ml-4 mt-2 rounded bg-green-500 px-4 py-2 text-white hover:bg-green-600"
            on:click={joinLobby}
        >
            Join Lobby
        </button>
    </div>
</main>
