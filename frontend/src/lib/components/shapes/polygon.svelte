<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let points: { x: number; y: number; isBeingModified: boolean }[];
	export let hexColor: string;
	export let rotationDegrees: number = 0;

	type CanvasContext = {
		registerDrawFunction: (fn: DrawFunction) => () => void;
		redrawCanvas: () => void;
	};

	const { registerDrawFunction, redrawCanvas } = getContext<CanvasContext>('canvas');

	$: if ((points, hexColor)) {
		redrawCanvas();
	}

	onMount(() => {
		const unregister = registerDrawFunction(draw);
		return () => {
			unregister();
			redrawCanvas();
		};
	});

	function getCentroid(pts: { x: number; y: number }[]): { x: number; y: number } {
		let centroid = { x: 0, y: 0 };
		pts.forEach((p) => {
			centroid.x += p.x;
			centroid.y += p.y;
		});
		centroid.x /= pts.length;
		centroid.y /= pts.length;
		return centroid;
	}

	function rotatePoint(
		point: { x: number; y: number },
		center: { x: number; y: number },
		angle: number
	) {
		const sin = Math.sin(angle);
		const cos = Math.cos(angle);
		const x = point.x - center.x;
		const y = point.y - center.y;
		const xNew = x * cos - y * sin;
		const yNew = x * sin + y * cos;
		return {
			x: xNew + center.x,
			y: yNew + center.y
		};
	}

	function draw(ctx: CanvasRenderingContext2D | null) {
		if (ctx) {
			ctx.fillStyle = hexColor;
			ctx.beginPath();

			// Check if points need rotation
			const c = getCentroid(points);
			const rotatedPoints = points.map((p) => rotatePoint(p, c, rotationDegrees));

			ctx.moveTo(rotatedPoints[0].x, rotatedPoints[0].y);
			for (let i = 1; i < rotatedPoints.length; i++) {
				const p = rotatedPoints[i];

				ctx.lineTo(p.x, p.y);
				ctx.fill();

				ctx.fillStyle = hexColor;
			}
			ctx.closePath();
			ctx.fill();

			// DOTS
			ctx.fillStyle = 'blue';
			for (let i = 0; i < rotatedPoints.length; i++) {
				const p = rotatedPoints[i];
				ctx.beginPath();
				ctx.arc(p.x, p.y, 5, 0, 2 * Math.PI); // Start point
				ctx.fill();
				ctx.closePath();
			}
		}
	}
</script>
