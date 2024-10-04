<script lang="ts">
	let actionName: string = 'Triangle';
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
		><path d="m80-160 401-640 399 640H80Zm107-60h586L481-685 187-220Zm293-233Z" /></svg
	>
</div>
