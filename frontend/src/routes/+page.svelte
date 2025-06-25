<script lang="ts">
    import { connectWebSocket } from '$lib/websocket';

    let username = '';
    let lobbyID = '';
</script>

<main class="flex h-screen flex-col items-center justify-center bg-gray-800 text-white">
    <h1 class="mb-4 text-center text-3xl font-bold">Welcome to WikiWoyage!</h1>

    <div class="mb-4">
        <form
            on:submit|preventDefault={() => connectWebSocket(username, lobbyID)}
            class="flex flex-col items-center gap-2"
        >
            <input
                type="text"
                id="username"
                bind:value={username}
                class="mt-2 w-64 rounded border border-gray-300 p-2"
                placeholder="Enter username"
            />
            <input
                type="text"
                id="lobbyid"
                bind:value={lobbyID}
                class="mt-2 w-64 rounded border border-gray-300 p-2"
                placeholder="Enter Lobby ID (optional)"
            />
            <button
                type="submit"
                class="mt-2 w-64 rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
            >
                {lobbyID ? 'Join Lobby' : 'Create Lobby'}
            </button>
        </form>
    </div>
</main>
