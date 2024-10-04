<script lang="ts">
	let actionName: string = 'Move';
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
			d="M480-80 317-243l44-44 89 89v-189h60v189l89-89 44 44L480-80ZM238-322 80-480l159-159 44 44-85 85h189v60H198l84 84-44 44Zm484 0-44-44 84-84H574v-60h188l-84-84 44-44 158 158-158 158ZM450-574v-188l-84 84-44-44 158-158 158 158-44 44-84-84v188h-60Z"
		/></svg
	>
</div>
