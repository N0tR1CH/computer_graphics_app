<script lang="ts">
	let actionName: string = 'QuadraticCurve';
	export let activeAction: string;
</script>

<div
	class="h-12 hover:bg-blue-500 rounded-xl hover:animate-pulse"
	class:bg-blue-500={activeAction === actionName}
>
	<svg
		class="hover:fill-white"
		class:fill-white={activeAction === actionName}
		xmlns="http://www.w3.org/2000/svg"
		height="48px"
		viewBox="0 -960 960 960"
		width="48px"
		fill="#434343"
		><path
			d="M720-160q0-116-44-218T556-556q-76-76-178-120t-218-44v-80q132 0 248.5 50.5T612-612q87 87 137.5 203.5T800-160h-80Z"
		/></svg
	>
</div>
