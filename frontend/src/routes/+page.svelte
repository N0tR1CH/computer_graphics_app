<script lang="ts">
	import TopBar from '$lib/components/top_bar.svelte';
	import ToolBar from '$lib/components/tool_bar.svelte';
	import ToolBarButton from '$lib/components/tool_bar_button.svelte';
	import TriangleOutline from '$lib/components/outlines/triangle_outline.svelte';
	import RectangleOutline from '$lib/components/outlines/rectangle_outline.svelte';
	import CircleOutline from '$lib/components/outlines/circle_outline.svelte';
	import MoveOutline from '$lib/components/outlines/move_outline.svelte';
	import ResizeOutline from '$lib/components/outlines/resize_outline.svelte';
	import StraightLineOutline from '$lib/components/outlines/straight_line_outline.svelte';
	import PencilOutline from '$lib/components/outlines/pencil_outline.svelte';
	import SaveOutline from '$lib/components/outlines/save_outline.svelte';
	import TextOutline from '$lib/components/outlines/text_outline.svelte';
	import Canvas from '$lib/components/canvas.svelte';
	import Rectangle from '$lib/components/shapes/rectangle.svelte';
	import Triangle from '$lib/components/shapes/triangle.svelte';
	import Ellipse from '$lib/components/shapes/ellipse.svelte';
	import StraightLine from '$lib/components/shapes/straight_line.svelte';
	import Text from '$lib/components/shapes/text.svelte';
	import type { Shape } from '../types/shape';
	import type { PossibleActions } from '../types/possible_actions';
	import ColorPickers from '$lib/components/color_picking/color_pickers.svelte';
	import ThirdDimensionCanvas from '$lib/components/third_dimension_canvas.svelte';
	import Image from '$lib/components/shapes/image.svelte';

	let activeAction: PossibleActions = 'Triangle';
	let text: string = '';
	let shapes: Shape[] = [];
	let sceneWidth: number = 350;
	let sceneHeight: number = 350;
</script>

<TopBar>
	<p class="text-2xl" slot="title">Stuff ;)</p>
</TopBar>
<ToolBar>
	<ToolBarButton bind:activeAction action={'Triangle'}>
		<TriangleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Rectangle'}>
		<RectangleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Circle'}>
		<CircleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Move'}>
		<MoveOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Resize'}>
		<ResizeOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'StraightLine'}>
		<StraightLineOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Pencil'}>
		<PencilOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Save'}>
		<SaveOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Text'}>
		<TextOutline {activeAction} />
	</ToolBarButton>
</ToolBar>
<Canvas height={500} width={1240} bind:shapes bind:activeAction bind:text>
	{#each shapes as shape}
		{#if shape.name === 'Rectangle'}
			<Rectangle
				x={shape.x}
				y={shape.y}
				height={shape.height}
				width={shape.width}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'Triangle'}
			<Triangle
				x={shape.x}
				y={shape.y}
				base={shape.base}
				height={shape.height}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'Ellipse'}
			<Ellipse
				x={shape.x}
				y={shape.y}
				radius1={shape.radius1}
				radius2={shape.radius2}
				rotation={shape.rotation}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'StraightLine'}
			<StraightLine x={shape.x} y={shape.y} x1={shape.x1} y1={shape.y1} hexColor={shape.hexColor} />
		{:else if shape.name === 'Text'}
			<Text x={shape.x} y={shape.y} text={shape.text} hexColor={shape.hexColor} />
		{:else if shape.name === 'Image'}
			<Image x={shape.x} y={shape.y} baseUrlImage={shape.baseUrlImage} />
		{/if}
	{/each}
</Canvas>
{#if activeAction === 'Move'}
	<p class="mt-4 text-center text-white">
		Current mode is moving elements. Press left click in order to move latest element somewhere else
		on the canvas.
	</p>
{:else if activeAction === 'Resize'}
	<p class="mt-4 text-center text-white">
		Current mode is resizing elements. Press left click in order to scale latest element. If it
		doesn't do anything that means element is not scalable!
	</p>
{:else if activeAction === 'Save'}
	<p class="mt-4 text-center text-white">
		Current mode is saving canvas. Click on the canvas in order to invoke action!
	</p>
{:else if activeAction === 'Text'}
	<p class="mt-4 text-center text-white">
		Write text to input field and click on the canvas in order to add it to correct coordinates.
	</p>
	<input
		class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
		bind:value={text}
		type="text"
	/>
{/if}
<ColorPickers></ColorPickers>
{#key sceneWidth + sceneHeight}
	<ThirdDimensionCanvas bind:shapes width={sceneWidth} height={sceneHeight}></ThirdDimensionCanvas>
{/key}
<div class="flex justify-center gap-x-4 mb-4">
	<label for="" class="text-white font-bold">
		Width of the scene
		<input type="number" bind:value={sceneWidth} class="rounded-xl p-2 text-black" />
	</label>
	<label for="" class="text-white font-bold">
		Height of the scene
		<input type="number" bind:value={sceneHeight} class="rounded-xl p-2 text-black" />
	</label>
</div>
<p class="text-center text-white mb-4 font-bold">
	Cube is being added to fixed position but you can change it by using move tool although it is
	clunky (3d to 2d ehh)
</p>

<style lang="postcss">
	:global(html) {
		background-color: rgba(0, 0, 0, 0);
	}
</style>
