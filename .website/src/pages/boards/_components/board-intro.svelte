<script lang="ts">
	import type { Snippet } from 'svelte';
	import BoardSection from '../../../components/board-section/board-section.svelte';
	type Props = {
		children: Snippet;
	};

	let { children }: Props = $props();
</script>

<BoardSection title="About this board" columnLabels={null}>
	<div class="thread-body board-intro">{@render children()}</div>
</BoardSection>

<style src="./board-intro.css"></style>
