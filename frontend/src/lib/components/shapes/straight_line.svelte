<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function';

	export let x: number;
	export let y: number;
	export let x1: number;
	export let y1: number;
	export let hexColor: string;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};
	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((x, y, x1, y1)) {
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
			if (x === x1 && y === y1) {
				// Draw a dot (small circle)
				const radius = 2; // Adjust the radius as needed
				ctx.beginPath();
				ctx.arc(x, y, radius, 0, 2 * Math.PI);
				ctx.fillStyle = hexColor;
				ctx.fill(); // Use fill to make it a solid dot
			} else {
				// Draw the line as usual
				ctx.beginPath();
				ctx.moveTo(x, y);
				ctx.lineTo(x1, y1);
				ctx.closePath();
				ctx.strokeStyle = hexColor;
				ctx.stroke();
			}
		}
	}
</script>
