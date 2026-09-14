<script lang="ts">
	import PixelFolder from '../pixel-folder/pixel-folder.svelte';

	type Props = {
		index: number;
		decisionDate: string | null;
		metadata: string | null;
		title: string;
		href: string;
		updated: boolean;
		lastUpdatedDate: string;
		onVisit: () => void;
	};

	let { index, title, href, updated, lastUpdatedDate, decisionDate, metadata, onVisit }: Props =
		$props();
</script>

<div class="thread-row" data-stripe={index % 2 === 0 ? 'a' : 'b'}>
	<PixelFolder {updated} pinned={false} width={19} />
	<a class="thread-row-title" {href} onclick={onVisit}>{title}</a>
	<span class="thread-row-status">{metadata ?? ''}</span>
	<span class="thread-row-date">
		{#if decisionDate !== null}decided {decisionDate}{:else}updated {lastUpdatedDate}{/if}
	</span>
</div>

<style src="./thread-row.css"></style>
