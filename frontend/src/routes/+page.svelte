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

	let activeAction: PossibleActions = 'Triangle';
	let text: string = '';

	let shapes: Shape[] = [];
</script>

<TopBar title={'Stuff ;)'} />
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
<Canvas height={500} width={1240} bind:shapes bind:activeAction>
	{#each shapes as shape}
		{#if shape.name === 'Rectangle'}
			<Rectangle x={shape.x} y={shape.y} height={shape.height} width={shape.width} />
		{:else if shape.name === 'Triangle'}
			<Triangle x={shape.x} y={shape.y} base={shape.base} height={shape.height} />
		{:else if shape.name === 'Ellipse'}
			<Ellipse
				x={shape.x}
				y={shape.y}
				radius1={shape.radius1}
				radius2={shape.radius2}
				rotation={shape.rotation}
			/>
		{:else if shape.name === 'StraightLine'}
			<StraightLine x={shape.x} y={shape.y} x1={shape.x1} y1={shape.y1} />
		{/if}
	{/each}
</Canvas>
{#if activeAction === 'Move'}
	<p class="text-white text-center mt-4">
		Current mode is moving elements. Press left click in order to move latest element somewhere else
		on the canvas.
	</p>
{:else if activeAction === 'Resize'}
	<p class="text-white text-center mt-4">
		Current mode is resizing elements. Press left click in order to scale latest element. If it
		doesn't do anything that means element is not scalable!
	</p>
{:else if activeAction === 'Save'}
	<p class="text-white text-center mt-4">
		Current mode is saving canvas. Click on the canvas in order to invoke action!
	</p>
{:else if activeAction === 'Text'}
	<p class="text-white text-center mt-4">
		Write text to input field and click on the canvas in order to add it to correct coordinates.
	</p>
	<input
		class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
		bind:value={text}
		type="text"
	/>
{/if}

<style lang="postcss">
	:global(html) {
		background-color: rgba(0, 0, 0, 0);
	}
</style>
