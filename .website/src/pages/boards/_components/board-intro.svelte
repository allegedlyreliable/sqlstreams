<script lang="ts">
	import type { Snippet } from 'svelte';
	import BoardSection from '../../../components/board-section/board-section.svelte';
	type Props = {
		title: string;
		children: Snippet;
	};

	let { title, children }: Props = $props();
</script>

<BoardSection {title} columnLabels={null}>
	<div class="thread-body board-intro">{@render children()}</div>
</BoardSection>

<style src="./board-intro.css"></style>
