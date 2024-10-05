<script lang="ts">
	import { setContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../types/draw_function';
	export let width = 0;
	export let height = 0;
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;
	let scheduled = false;

	const drawFunctions = new Set<DrawFunction>();
	function registerDrawFunction(fn: DrawFunction) {
		drawFunctions.add(fn);
		return () => {
			drawFunctions.delete(fn);
		};
	}

	function redrawCanvas() {
		if (scheduled) return;
		scheduled = true;
		requestAnimationFrame(() => {
			draw();
			scheduled = false;
		});
	}

	setContext('canvas', {
		registerDrawFunction,
		redrawCanvas
	});

	function draw() {
		if (ctx) {
			ctx.clearRect(0, 0, canvas.width, canvas.height);
			drawFunctions.forEach((fn) => fn(ctx));
		}
	}

	$: if (height > 500) {
		draw();
	}

	onMount(() => {
		canvas.style.width = '100%';
		canvas.style.height = '100%';
		canvas.width = canvas.offsetWidth;
		canvas.height = canvas.offsetHeight;

		ctx = canvas.getContext('2d');
		draw();
	});
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
