import { writable } from 'svelte/store';

export const currentColor = writable<string>('#000000');
