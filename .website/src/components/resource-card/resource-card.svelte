<script lang="ts">
	import { cardGeometry } from './geometry';
	import type { ResourceNode } from '../diagram/types';

	type Props = { data: ResourceNode['data'] };
	let { data }: Props = $props();
</script>

<div
	class="resource"
	style:--header-height={`${cardGeometry.headerHeight}px`}
	style:--line-height={`${cardGeometry.lineHeight}px`}
	style:--body-padding={`${cardGeometry.bodyPadding}px`}
	style:--border-width={`${cardGeometry.borderWidth}px`}
>
	<div class="header">{data.title}</div>
	<div class="body">
		{#if data.lines !== null}
			<div>
				{#each data.lines as line, index (index)}<div>{line}</div>{/each}
			</div>
		{/if}
	</div>
</div>

<style src="./resource-card.css"></style>
