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

	let oldPos = { x: 0, y: 0 };
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;
	let scheduled = false;
	let cursorPosition = { x: 0, y: 0 };
	let isDrawing = false;
	let isLive = false;
	let isPencil = false;

	/* Start of Moving bezier state */
	let isMovingStart: boolean = false;
	let isMovingEnd: boolean = false;
	let isMovingCp1: boolean = false;
	let isMovingCp2: boolean = false;
	/* End of Moving bezier state */

	function calculateControlPoints(
		start: { x: number; y: number },
		end: { x: number; y: number },
		offsetX: number,
		offsetY: number
	) {
		const cp1 = {
			x: start.x + offsetX,
			y: start.y + offsetY
		};

		const cp2 = {
			x: end.x - offsetX,
			y: end.y - offsetY
		};

		return { cp1, cp2 };
	}

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
				case 'Bezier':
					newShapeName = 'Bezier';
					break;
				case 'QuadraticCurve':
					newShapeName = 'QuadraticCurve';
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
					baseUrlImage: '',
					bezierStart: { x: cursorPosition.x, y: cursorPosition.y },
					bezierCp1: { x: 0, y: 0 },
					bezierCp2: { x: 0, y: 0 },
					bezierEnd: { x: 0, y: 0 }
				}
			];

			if (activeAction === 'Bezier') {
				const endPoint = { x: cursorPosition.x + 100, y: cursorPosition.y + 100 }; // Example end point
				const offsetX = 50;
				const offsetY = 50;
				const controlPoints = calculateControlPoints(
					shapes[shapes.length - 1].bezierStart,
					endPoint,
					offsetX,
					offsetY
				);
				shapes[shapes.length - 1].bezierCp1 = controlPoints.cp1;
				shapes[shapes.length - 1].bezierCp2 = controlPoints.cp2;
				shapes[shapes.length - 1].bezierEnd = endPoint;
			}

			if (activeAction === 'QuadraticCurve') {
				const endPoint = { x: cursorPosition.x + 100, y: cursorPosition.y + 100 }; // Example end point
				const startP = shapes[shapes.length - 1].bezierStart;
				let controlX = 0.5 * startP.x + 0.5 * endPoint.x;
				let controlY = 0.5 * startP.y + 0.5 * endPoint.y;

				shapes[shapes.length - 1].bezierCp1 = { x: controlX, y: controlY };
				shapes[shapes.length - 1].bezierEnd = endPoint;
			}
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
					x1: oldPos.x,
					y1: oldPos.y,
					text: '',
					hexColor: $currentColor,
					baseUrlImage: ''
				}
			];
			oldPos.x = cursorPosition.x;
			oldPos.y = cursorPosition.y;
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

		const bezierCurveAction = () => {
			const lastShape = shapes[shapes.length - 1];
			if (lastShape.name !== 'Bezier') return;

			const endPoint = { x: cursorPosition.x, y: cursorPosition.y };
			const offsetX = 50;
			const offsetY = 50;

			const controlPoints = calculateControlPoints(
				lastShape.bezierStart,
				endPoint,
				offsetX,
				offsetY
			);
			lastShape.bezierCp1 = controlPoints.cp1;
			lastShape.bezierCp2 = controlPoints.cp2;
			lastShape.bezierEnd = endPoint;
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
			case 'Bezier': {
				bezierCurveAction();
				break;
			}
			case 'Move': {
				shapes[shapes.length - 1].x = cursorPosition.x;
				shapes[shapes.length - 1].y = cursorPosition.y;

				if (shapes[shapes.length - 1].name == 'Bezier') {
					if (isMovingStart) {
						shapes[shapes.length - 1].bezierStart.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierStart.y = cursorPosition.y;
					} else if (isMovingEnd) {
						shapes[shapes.length - 1].bezierEnd.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierEnd.y = cursorPosition.y;
					} else if (isMovingCp1) {
						shapes[shapes.length - 1].bezierCp1.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierCp1.y = cursorPosition.y;
					} else if (isMovingCp2) {
						shapes[shapes.length - 1].bezierCp2.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierCp2.y = cursorPosition.y;
					}
				}

				if (shapes[shapes.length - 1].name == 'QuadraticCurve') {
					if (isMovingStart) {
						shapes[shapes.length - 1].bezierStart.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierStart.y = cursorPosition.y;
					} else if (isMovingEnd) {
						shapes[shapes.length - 1].bezierEnd.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierEnd.y = cursorPosition.y;
					} else if (isMovingCp1) {
						shapes[shapes.length - 1].bezierCp1.x = cursorPosition.x;
						shapes[shapes.length - 1].bezierCp1.y = cursorPosition.y;
					}
				}
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
			const s = shapes[shapes.length - 1];
			if (s.name === 'Bezier') {
				if (
					Math.abs(s.bezierStart.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierStart.y - cursorPosition.y) < 5
				) {
					isMovingStart = true;
				} else if (
					Math.abs(s.bezierEnd.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierEnd.y - cursorPosition.y) < 5
				) {
					isMovingEnd = true;
				} else if (
					Math.abs(s.bezierCp1.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierCp1.y - cursorPosition.y) < 5
				) {
					isMovingCp1 = true;
				} else if (
					Math.abs(s.bezierCp2.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierCp2.y - cursorPosition.y) < 5
				) {
					isMovingCp2 = true;
				}
			}

			if (s.name === 'QuadraticCurve') {
				if (
					Math.abs(s.bezierStart.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierStart.y - cursorPosition.y) < 5
				) {
					isMovingStart = true;
				} else if (
					Math.abs(s.bezierEnd.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierEnd.y - cursorPosition.y) < 5
				) {
					isMovingEnd = true;
				} else if (
					Math.abs(s.bezierCp1.x - cursorPosition.x) < 5 &&
					Math.abs(s.bezierCp1.y - cursorPosition.y) < 5
				) {
					isMovingCp1 = true;
				}
			}
			isLive = true;
			return;
		}

		if (activeAction === 'Pencil') {
			oldPos.x = cursorPosition.x;
			oldPos.y = cursorPosition.y;
			isPencil = true;
		}

		if (isDrawing) {
			return;
		}

		isDrawing = true;
		drawing();
	}}
	on:mouseup={() => {
		isMovingStart = false;
		isMovingEnd = false;
		isMovingCp1 = false;
		isMovingCp2 = false;

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
