<script lang="ts">
	import { slopCopy, type SlopNoticeLevel } from './types';

	type Props = {
		level: SlopNoticeLevel;
	};

	let { level }: Props = $props();

	const copy = $derived(slopCopy[level]);
</script>

<aside class="slop-notice" data-level={level} data-pagefind-ignore>
	<div class="notice-header">
		<span class="notice-chip">SLOP LEVEL</span>
		<span class="notice-title">{level}</span>
		<svg class="notice-meter" width="26" height="8" viewBox="0 0 26 8" aria-hidden="true">
			{#each [0, 1, 2] as block (block)}
				<rect x={block * 9} y="0" width="8" height="8" data-filled={block < copy.blocks} />
			{/each}
		</svg>
	</div>
	<p class="notice-body">{copy.body}</p>
</aside>

<style src="./slop-notice.css"></style>
