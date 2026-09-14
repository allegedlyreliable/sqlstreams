<script lang="ts">
	import PixelFolder from '../pixel-folder/pixel-folder.svelte';

	type Props = {
		docCount: number;
		codeCount: number;
	};

	let { docCount, codeCount }: Props = $props();
</script>

<div class="board-stats">
	<div class="stats">
		<span class="stats-line">
			<a href="/whats-new/">{docCount} documents</a>
		</span>
		{#if codeCount > 0}
			<span class="stats-line">
				<a href="/boards/troubleshooting/">{codeCount} diagnostic codes</a>
			</span>
		{/if}
	</div>
	<div class="legend">
		<span class="legend-row">
			<PixelFolder updated={true} pinned={false} width={13} /> updated since your last visit
		</span>
		<span class="legend-row">
			<PixelFolder updated={false} pinned={false} width={13} /> read
		</span>
		<span class="legend-row legend-visited">purple links — already visited</span>
	</div>
</div>

<style src="./board-stats.css"></style>
