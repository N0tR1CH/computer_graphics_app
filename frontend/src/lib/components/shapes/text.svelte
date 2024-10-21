<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function';

	export let x: number;
	export let y: number;
	export let text: string;
	export let hexColor: string;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};
	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((x, y, text)) {
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
			ctx.fillStyle = hexColor;
			ctx.fillText(text, x, y);
		}
	}
</script>
