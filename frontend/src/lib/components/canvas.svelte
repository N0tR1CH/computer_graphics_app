<script lang="ts">
	import { onMount, setContext } from 'svelte';

	export let width: `${number}` | number = 0;
	export let height: `${number}` | number = 0;
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;

	onMount(() => {
		ctx = canvas.getContext('2d');
	});

	setContext('canvas', () => ctx);
</script>

<canvas
	class="rounded-xl shadow-xl bg-white"
	class:w-full={width === 0}
	bind:this={canvas}
	{width}
	{height}
	on:click={() => {
		console.log('canvas clicked');
	}}
></canvas>
<slot />
