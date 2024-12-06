<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import type { DrawFunction } from '../../../types/draw_function.ts';

	export let points: { x: number; y: number; isBeingModified: boolean }[];
	export let hexColor: string;
	export let rotationDegrees: number = 0;
	export let scale: number = 2;
	export let offsetX: number = 0;
	export let offsetY: number = 0;

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
		isBeingModified: boolean,
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
			y: yNew + center.y,
			isBeingModified: isBeingModified
		};
	}

	function scalePoint(
		point: { x: number; y: number },
		center: { x: number; y: number },
		isBeingModified: boolean,
		scale: number
	) {
		const x = point.x - center.x;
		const y = point.y - center.y;
		const xScaled = x * scale;
		const yScaled = y * scale;

		return {
			x: xScaled + center.x,
			y: yScaled + center.y,
			isBeingModified: isBeingModified
		};
	}

	function draw(ctx: CanvasRenderingContext2D | null) {
		if (ctx) {
			ctx.fillStyle = hexColor;
			ctx.beginPath();

			const c = getCentroid(points);
			/* PIPELINE
			 *  rotation -> scaling -> offset
			 */
			const rotatedPoints = points.map((p) =>
				rotatePoint(p, c, p.isBeingModified, rotationDegrees)
			);
			const scaledPoints = rotatedPoints.map((p) => scalePoint(p, c, p.isBeingModified, scale));
			const offsetPoints = scaledPoints.map((p) => ({
				x: p.x + offsetX,
				y: p.y + offsetY,
				isBeingModified: p.isBeingModified
			}));

			ctx.moveTo(offsetPoints[0].x, offsetPoints[0].y);
			for (let i = 1; i < offsetPoints.length; i++) {
				const p = offsetPoints[i];

				ctx.lineTo(p.x, p.y);
				ctx.fill();

				ctx.fillStyle = hexColor;
			}
			ctx.closePath();
			ctx.fill();

			// DOTS
			ctx.fillStyle = 'blue';
			for (let i = 0; i < offsetPoints.length; i++) {
				const p = offsetPoints[i];
				ctx.beginPath();
				ctx.arc(p.x, p.y, 5, 0, 2 * Math.PI); // Start point
				ctx.fill();
				ctx.closePath();
			}
		}
	}
</script>
