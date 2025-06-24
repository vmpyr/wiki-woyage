<script lang="ts">
	import { onMount } from 'svelte';

	let username = '';
	let lobbyID = '';

	onMount(() => {
		username = localStorage.getItem('username') || '';
		lobbyID = localStorage.getItem('lobbyID') || '';

		if (!username || !lobbyID) {
			window.location.href = '/';
		}
	});
</script>

<main class="flex h-screen flex-col items-center justify-center bg-gray-800 text-white">
	<h1 class="mb-4 text-center text-3xl font-bold">Welcome {username} to Lobby {lobbyID}!</h1>
</main>
