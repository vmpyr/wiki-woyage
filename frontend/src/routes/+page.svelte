<script>
    import { onMount } from 'svelte';
    import { connectWebSocket, sendWebSocketMessage } from '$lib/websocket';

    let playerName = '';

    onMount(() => {
        connectWebSocket();
    });

    function createPlayer() {
        sendWebSocketMessage('create_player', { playerName });
    }
</script>

<main class="flex h-screen flex-col items-center justify-center bg-gray-800 text-white">
    <h1 class="mb-4 text-3xl font-bold">Welcome to WikiWoyage!</h1>

    <div class="mb-4">
        <form on:submit|preventDefault={createPlayer} class="flex flex-row items-center gap-2">
            <input
                type="text"
                id="playerName"
                bind:value={playerName}
                class="mt-2 w-64 rounded border border-gray-300 p-2"
                placeholder="Enter username"
            />
            <button
                type="submit"
                class="mt-2 rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
            >
                Let's go!
            </button>
        </form>
    </div>
</main>
