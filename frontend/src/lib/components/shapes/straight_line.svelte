<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function';

	export let x: number;
	export let y: number;
	export let x1: number;
	export let y1: number;

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
			ctx.beginPath();
			ctx.moveTo(x, y);
			ctx.lineTo(x1, y1);
			ctx.closePath();
			ctx.stroke();
		}
	}
</script>
