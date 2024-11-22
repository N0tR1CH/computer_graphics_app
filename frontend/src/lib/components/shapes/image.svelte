<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let baseUrlImage: string;
	export let x: number = 0;
	export let y: number = 0;

	type CanvasContextType = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};

	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContextType>('canvas');

	let image: HTMLImageElement | null = null;
	let unregister: (() => void) | null = null;

	onMount(() => {
		image = new Image();
		image.src = baseUrlImage;

		image.onload = () => {
			// Register the draw function after the image has loaded
			registerDraw();
		};

		image.onerror = (error) => {
			console.error('Error loading image:', error);
		};

		return () => {
			if (unregister) {
				unregister();
			}
			redrawCanvas();
		};
	});

	$: if (image && baseUrlImage) {
		image.src = baseUrlImage;
		image.onload = () => {
			registerDraw();
		};
	}

	$: if ((x || y) && image) {
		if (image.complete) {
			registerDraw();
		}
	}

	function registerDraw() {
		unregister = registerDrawFunction(draw);
		redrawCanvas(); // Trigger a redraw after registering the draw function
	}

	function draw(ctx: CanvasRenderingContext2D | null) {
		if (ctx && image && image.complete) {
			ctx.drawImage(image, x, y);
		}
	}
</script>
