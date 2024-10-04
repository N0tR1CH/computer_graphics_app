<script lang="ts">
	let actionName: string = 'Text';
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
			d="M290-160v-540H80v-100h520v100H390v540H290Zm360 0v-340H520v-100h360v100H750v340H650Z"
		/></svg
	>
</div>
