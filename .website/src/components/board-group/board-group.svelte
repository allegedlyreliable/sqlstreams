<script lang="ts">
	import ThreadRows from '../thread-rows/thread-rows.svelte';
	import type { BoardGroupSection } from './types';
	type Props = { id: string; title: string; sections: BoardGroupSection[] };
	let { id, title, sections }: Props = $props();
</script>

<section class="board-group" {id} aria-labelledby={`${id}-title`}>
	<header class="group-heading">
		<h2 id={`${id}-title`}>{title}</h2>
		<span>Updated</span>
	</header>
	{#each sections as section (section.id)}
		<div id={section.id ?? undefined}>
			{#if section.title !== null}
				<header class="section-heading">
					<h3>{section.title}</h3>
				</header>
			{/if}
			<ThreadRows rows={section.rows} />
		</div>
	{/each}
</section>

<style src="./board-group.css"></style>
