<script lang="ts">
	export let action: string;
	export let activeAction: string;
</script>

<li>
	<button
		on:click={() => {
			activeAction = action;
		}}
	>
		<slot></slot>
	</button>
</li>
