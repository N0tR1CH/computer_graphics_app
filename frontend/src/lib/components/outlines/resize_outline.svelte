<script lang="ts">
	let actionName: string = 'Resize';
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
		><path d="M120-120v-300h60v198l558-558H540v-60h300v300h-60v-198L222-180h198v60H120Z" /></svg
	>
</div>
