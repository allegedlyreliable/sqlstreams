<script lang="ts">
	import GroupMapSection from '../group-map-section/group-map-section.svelte';
	import type { GroupMapData } from './types';
	type Props = { map: GroupMapData; currentHref: string | null };
	let { map, currentHref }: Props = $props();
</script>

<details class="group-map" data-pagefind-ignore>
	<summary>{map.title} · Related pages</summary>
	<nav aria-label={`${map.title} related pages`}>
		{#each map.sections as section (section.title)}
			<GroupMapSection {section} {currentHref} />
		{/each}
	</nav>
</details>

<style src="./group-map.css"></style>
