<script lang="ts">
    import { onMount } from 'svelte';
    import { getClientInfo } from '$lib/lobby';

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
</main>
