<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let height: number;
	export let width: number;
	export let x: number;
	export let y: number;
	export let fill: string;
	export let stroke: string;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};

	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((height, width, x, y, fill, stroke)) {
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
			ctx.fillStyle = fill;
			ctx.strokeStyle = stroke;
			ctx.beginPath();
			ctx.rect(x, y, height, width);
			ctx.stroke();
			ctx.fill();
		}
	}
</script>
