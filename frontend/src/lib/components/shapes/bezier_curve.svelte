<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let start: { x: number; y: number };
	export let cp1: { x: number; y: number };
	export let cp2: { x: number; y: number };
	export let end: { x: number; y: number };
	export let hexColor: string;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};

	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((start, cp1, cp2, end)) {
		redrawCanvas();
	}

	onMount(() => {
		const unregister = registerDrawFunction(draw);
		return () => {
			unregister();
			redrawCanvas();
		};
	});

	function draw(ctx: CanvasRenderingContext2D | null) {
		if (ctx) {
			// src: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/bezierCurveTo

			// Cubic Bézier curve
			ctx.fillStyle = hexColor;
			ctx.beginPath();
			ctx.moveTo(start.x, start.y);
			ctx.bezierCurveTo(cp1.x, cp1.y, cp2.x, cp2.y, end.x, end.y);
			ctx.strokeStyle = hexColor;
			ctx.stroke();

			// Start and end points
			ctx.fillStyle = 'blue';
			ctx.beginPath();
			ctx.arc(start.x, start.y, 5, 0, 2 * Math.PI); // Start point
			ctx.arc(end.x, end.y, 5, 0, 2 * Math.PI); // End point
			ctx.fill();

			// Control points
			ctx.fillStyle = 'red';
			ctx.beginPath();
			ctx.arc(cp1.x, cp1.y, 5, 0, 2 * Math.PI); // Control point one
			ctx.arc(cp2.x, cp2.y, 5, 0, 2 * Math.PI); // Control point two
			ctx.fill();
		}
	}
</script>
