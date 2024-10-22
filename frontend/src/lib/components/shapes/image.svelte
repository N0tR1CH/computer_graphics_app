<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let baseUrlImage: string;
	export let x: number;
	export let y: number;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};

	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((x, y)) {
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
			const image = new Image();
			image.onload = () => {
				ctx.drawImage(image, x, y);
			};
			image.src = baseUrlImage;
		}
	}
</script>
