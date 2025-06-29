import { writable } from 'svelte/store';

export const playerList = writable<string[][]>([]);
