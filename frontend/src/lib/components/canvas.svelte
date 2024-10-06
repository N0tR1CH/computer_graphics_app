<script lang="ts">
	import { setContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../types/draw_function';
	import type { PossibleActions } from '../../types/possible_actions';
	import type { Shape } from '../../types/shape';

	export let width = 0;
	export let height = 0;
	export let shapes: Shape[];
	export let activeAction: PossibleActions;
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;
	let scheduled = false;
	let cursorPosition = { x: 0, y: 0 };
	let isDrawing = false;
	let isLive = false;

	function drawing() {
		if (isDrawing && !isLive) {
			isLive = true;
			let newShapeName: string;

			switch (activeAction) {
				case 'Triangle':
					newShapeName = 'Triangle';
					break;
				case 'Rectangle':
					newShapeName = 'Rectangle';
					break;
				case 'Circle':
					newShapeName = 'Ellipse';
					break;
				default:
					newShapeName = '';
			}

			if (newShapeName == '') {
				return;
			}

			shapes = [
				...shapes,
				{
					name: newShapeName,
					x: cursorPosition.x,
					y: cursorPosition.y,
					height: 0,
					width: 0,
					base: 0,
					radius1: 0,
					radius2: 0,
					rotation: 0
				}
			];
		} else {
			isLive = false;
			shapes = [...shapes];
		}
	}

	function handleMove(event: MouseEvent) {
		const rect = canvas.getBoundingClientRect();

		cursorPosition.x = event.clientX - rect.left;
		cursorPosition.y = event.clientY - rect.top;

		if (!isLive) {
			return;
		}

		const lastShape = shapes[shapes.length - 1];
		switch (activeAction) {
			case 'Rectangle': {
				const width = Math.abs(lastShape.x - cursorPosition.x);
				const height = Math.abs(lastShape.y - cursorPosition.y);

				if (shapes[shapes.length - 1].y - cursorPosition.y > 0) {
					shapes[shapes.length - 1].height = -height;
				} else {
					shapes[shapes.length - 1].height = height;
				}

				if (lastShape.x - cursorPosition.x > 0) {
					shapes[shapes.length - 1].width = -width;
				} else {
					shapes[shapes.length - 1].width = width;
				}
				break;
			}
			case 'Triangle': {
				const base = Math.abs(lastShape.x - cursorPosition.x);
				const height = Math.abs(lastShape.y - cursorPosition.y);
				if (shapes[shapes.length - 1].y - cursorPosition.y > 0) {
					shapes[shapes.length - 1].height = -height;
				} else {
					shapes[shapes.length - 1].height = height;
				}

				shapes[shapes.length - 1].base = base;
				break;
			}
			case 'Circle': {
				const radiusX = Math.abs(lastShape.x - cursorPosition.x);
				const radiusY = Math.abs(lastShape.y - cursorPosition.y);

				shapes[shapes.length - 1].radius1 = radiusX;
				shapes[shapes.length - 1].radius2 = radiusY;
				shapes[shapes.length - 1].rotation = Math.atan2(
					cursorPosition.y - lastShape.y,
					cursorPosition.x - lastShape.x
				);
			}
		}

		shapes = [...shapes];
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
	on:mousedown={() => {
		console.log('Mouse pressed');
		if (!isDrawing) {
			isDrawing = true;
			drawing();
		}
	}}
	on:mouseup={() => {
		console.log('Mouse released');
		if (isDrawing) {
			isDrawing = false;
			drawing();
		}
	}}
></canvas>
<slot />
