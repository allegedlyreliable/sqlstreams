<script lang="ts">
	import type { GroupMapSection } from '../group-map/types';

	type Props = {
		section: GroupMapSection;
		currentHref: string | null;
	};

	let { section, currentHref }: Props = $props();
	const containsCurrent = $derived(section.entries.some((entry) => entry.href === currentHref));
</script>

<details class="map-section" open={containsCurrent}>
	<summary>{section.title}</summary>
	<ul>
		{#each section.entries as entry (entry.href)}
			<li>
				<a href={entry.href} aria-current={entry.href === currentHref ? 'page' : undefined}>
					{entry.title}
				</a>
			</li>
		{/each}
	</ul>
</details>

<style src="./group-map-section.css"></style>
