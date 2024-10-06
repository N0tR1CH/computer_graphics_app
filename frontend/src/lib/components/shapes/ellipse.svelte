<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function';

	export let radius1: number;
	export let radius2: number;
	export let x: number;
	export let y: number;
	export let rotation: number = 0;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};
	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((radius1, radius2, rotation, x, y)) {
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
			ctx.ellipse(x, y, radius1, radius2, (rotation * Math.PI) / 180, 0, 2 * Math.PI);
			ctx.fill();
		}
	}
</script>
