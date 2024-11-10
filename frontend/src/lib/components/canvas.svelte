<script lang="ts">
	import type { DrawFunction } from '../../types/draw_function';
	import type { PossibleActions } from '../../types/possible_actions';
	import type { Shape } from '../../types/shape';
	import { currentColor } from '$lib/stores/stores';
	import { setContext, onMount } from 'svelte';
	import { SaveCanvasImg } from '$lib/wailsjs/go/main/App';
	import { main } from '$lib/wailsjs/go/models';

	export let text;
	export let width = 0;
	export let height = 0;
	export let shapes: Shape[];
	export let activeAction: PossibleActions;
	export let selectedFileFormat: main.ImageFormat;
	export let comments: string[] = [];

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;
	let scheduled = false;
	let cursorPosition = { x: 0, y: 0 };
	let isDrawing = false;
	let isLive = false;
	let isPencil = false;

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
				case 'StraightLine':
					newShapeName = 'StraightLine';
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
					rotation: 0,
					x1: 0,
					y1: 0,
					text: '',
					hexColor: $currentColor,
					baseUrlImage: ''
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

		if (activeAction === 'Pencil' && isPencil) {
			shapes = [
				...shapes,
				{
					name: 'StraightLine',
					x: cursorPosition.x,
					y: cursorPosition.y,
					height: 0,
					width: 0,
					base: 0,
					radius1: 0,
					radius2: 0,
					rotation: 0,
					x1: cursorPosition.x,
					y1: cursorPosition.y,
					text: '',
					hexColor: $currentColor,
					baseUrlImage: ''
				}
			];
			return;
		}

		if (!isLive) {
			return;
		}

		const lastShape = shapes[shapes.length - 1];
		const rectangleAction = () => {
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
		};

		const triangleAction = () => {
			const base = Math.abs(lastShape.x - cursorPosition.x);
			const height = Math.abs(lastShape.y - cursorPosition.y);
			if (shapes[shapes.length - 1].y - cursorPosition.y > 0) {
				shapes[shapes.length - 1].height = -height;
			} else {
				shapes[shapes.length - 1].height = height;
			}

			shapes[shapes.length - 1].base = base;
		};

		const circleAction = () => {
			const radiusX = Math.abs(lastShape.x - cursorPosition.x);
			const radiusY = Math.abs(lastShape.y - cursorPosition.y);

			shapes[shapes.length - 1].radius1 = radiusX;
			shapes[shapes.length - 1].radius2 = radiusY;
			shapes[shapes.length - 1].rotation = Math.atan2(
				cursorPosition.y - lastShape.y,
				cursorPosition.x - lastShape.x
			);
		};

		switch (activeAction) {
			case 'Rectangle': {
				rectangleAction();
				break;
			}
			case 'Triangle': {
				triangleAction();
				break;
			}
			case 'Circle': {
				circleAction();
				break;
			}
			case 'Move': {
				if (cursorPosition.x) shapes[shapes.length - 1].x = cursorPosition.x;
				shapes[shapes.length - 1].y = cursorPosition.y;
				break;
			}
			case 'Resize': {
				switch (lastShape.name) {
					case 'Triangle': {
						triangleAction();
						break;
					}
					case 'Rectangle': {
						rectangleAction();
						break;
					}
					case 'Ellipse': {
						circleAction();
						break;
					}
				}
				break;
			}
			case 'StraightLine': {
				shapes[shapes.length - 1].x1 = cursorPosition.x;
				shapes[shapes.length - 1].y1 = cursorPosition.y;
				break;
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
			ctx.fillStyle = 'white';
			ctx.fillRect(0, 0, canvas.width, canvas.height);
			ctx.fillStyle = 'black';
			drawFunctions.forEach((fn) => fn(ctx));
		}
	}

	$: if (height > 500) {
		draw();
	}

	onMount(() => {
		canvas.style.width = '100%';
		canvas.width = canvas.offsetWidth;
		ctx = canvas.getContext('2d');
		if (ctx) {
			draw();
			ctx.fillStyle = 'black';
		}
	});
</script>

<canvas
	class="rounded-xl bg-white shadow-xl"
	class:w-full={width === 0}
	bind:this={canvas}
	{width}
	{height}
	on:pointermove={handleMove}
	on:mousedown={() => {
		if (activeAction === 'Text') {
			shapes = [
				...shapes,
				{
					name: 'Text',
					x: cursorPosition.x,
					y: cursorPosition.y,
					height: 0,
					width: 0,
					base: 0,
					radius1: 0,
					radius2: 0,
					rotation: 0,
					x1: 0,
					y1: 0,
					text: text,
					hexColor: $currentColor,
					baseUrlImage: ''
				}
			];
			shapes = [...shapes];
			return;
		}

		if (activeAction === 'Save') {
			const dataURI = canvas.toDataURL('image/jpeg');
			SaveCanvasImg(dataURI, selectedFileFormat, comments);
			comments = [];
			return;
		}

		console.log('Mouse pressed');
		if (activeAction === 'Move' || activeAction === 'Resize') {
			isLive = true;
			return;
		}

		if (activeAction === 'Pencil') {
			isPencil = true;
		}

		if (isDrawing) {
			return;
		}

		isDrawing = true;
		drawing();
	}}
	on:mouseup={() => {
		if (activeAction === 'Text') {
			return;
		}
		console.log('Mouse released');
		if (activeAction === 'Move' || activeAction === 'Resize') {
			isLive = false;
		}

		if (!isDrawing) {
			return;
		}

		if (activeAction === 'Pencil') {
			isPencil = false;
			return;
		}

		isDrawing = false;
		if (['Triangle', 'Rectangle', 'Ellipse'].includes(shapes[shapes.length - 1].name)) {
			activeAction = 'Move';
		}
		drawing();
	}}
></canvas>
<slot />
