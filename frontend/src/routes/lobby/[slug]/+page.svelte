<script lang="ts">
    import { onMount } from 'svelte';
    import { getClientInfo, requestAdminStatus } from '$lib/lobby';
    import { playerList } from '$lib/stores';

    let username = '';
    let lobbyID = '';

    onMount(async () => {
        const clientInfo = await getClientInfo();
        if (!clientInfo) {
            window.location.href = '/';
            return;
        }
        username = clientInfo.username;
        lobbyID = clientInfo.lobbyID;
    });
</script>

<main class="flex h-screen flex-col items-center justify-center bg-gray-800 text-white">
    <h1 class="mb-4 text-center text-3xl font-bold">Welcome {username} to Lobby {lobbyID}!</h1>
    <div class="absolute left-8 top-1/2 w-60 -translate-y-1/2 rounded-lg bg-gray-700 p-4">
        <h2 class="mb-2 text-xl font-semibold">Players in Lobby:</h2>
        <ul>
            {#each $playerList as player}
                <li>{player} joined</li>
            {/each}
        </ul>
    </div>
    <button
        class="mt-8 rounded bg-blue-600 px-4 py-2 font-semibold hover:bg-blue-700"
        on:click={() => requestAdminStatus()}
    >
        Am I the admin?
    </button>
</main>
