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

	type PossibleActions =
		| 'Triangle'
		| 'Rectangle'
		| 'Circle'
		| 'Move'
		| 'Resize'
		| 'StraightLine'
		| 'Pencil'
		| 'Text';

	let activeAction: PossibleActions = 'Triangle';
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
<Canvas height={600}></Canvas>

<style lang="postcss">
	:global(html) {
		background-color: rgba(0, 0, 0, 0);
	}
</style>
