<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function';

	export let x: number;
	export let y: number;
	export let base: number;
	export let height: number;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};
	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((x, y, base, height)) {
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
			ctx.moveTo(x - base / 2, y + height / 2);
			ctx.lineTo(x, y - height / 2);
			ctx.lineTo(x + base / 2, y + height / 2);
			ctx.closePath();
			ctx.fill();
		}
	}
</script>
