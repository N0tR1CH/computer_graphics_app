<script lang="ts">
	import { setContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../types/draw_function';

	type Shape = {
		name: string;
		x: number;
		y: number;
		height: number;
		width: number;
	};

	export let width = 0;
	export let height = 0;
	export let shapes: Shape[];
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;
	let scheduled = false;
	let cursorPosition = { x: 0, y: 0 };
	let isDrawing = false;
	let isLive = false;

	function drawing(event: MouseEvent) {
		if (isDrawing) {
			console.log('I AM DRAWING');
			isLive = true;
			shapes = [
				...shapes,
				{ name: 'Rectangle', x: cursorPosition.x, y: cursorPosition.y, height: 0, width: 0 }
			];
		} else {
			console.log('I STOPPED DRAWING');
			const lastShape = shapes[shapes.length - 1];

			// Calculate width and height based on the current cursor position
			const width = Math.abs(lastShape.x - cursorPosition.x);
			const height = Math.abs(lastShape.y - cursorPosition.y);

			if (lastShape.y - cursorPosition.y > 0) {
				shapes[shapes.length - 1].height = -height;
			} else {
				shapes[shapes.length - 1].height = height;
			}

			if (lastShape.x - cursorPosition.x > 0) {
				shapes[shapes.length - 1].width = -width;
			} else {
				shapes[shapes.length - 1].width = width;
			}

			isLive = false;
			shapes = [...shapes];
		}
	}

	function handleMove(event: MouseEvent) {
		const rect = canvas.getBoundingClientRect();

		cursorPosition.x = event.clientX - rect.left;
		cursorPosition.y = event.clientY - rect.top;

		if (isLive) {
			const lastShape = shapes[shapes.length - 1];

			// Calculate width and height based on the current cursor position
			const width = Math.abs(lastShape.x - cursorPosition.x);
			const height = Math.abs(lastShape.y - cursorPosition.y);

			if (lastShape.y - cursorPosition.y > 0) {
				shapes[shapes.length - 1].height = -height;
			} else {
				shapes[shapes.length - 1].height = height;
			}

			if (lastShape.x - cursorPosition.x > 0) {
				shapes[shapes.length - 1].width = -width;
			} else {
				shapes[shapes.length - 1].width = width;
			}

			shapes = [...shapes];
		}
	}

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
	on:pointermove={handleMove}
	on:mousedown={(e) => {
		console.log('Mouse pressed');
		if (!isDrawing) {
			isDrawing = true;
			drawing(e);
		}
	}}
	on:mouseup={(e) => {
		console.log('Mouse released');
		if (isDrawing) {
			isDrawing = false;
			drawing(e);
		}
	}}
></canvas>
<slot />
