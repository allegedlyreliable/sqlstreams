<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import BoardStats from './board-stats.svelte';

	const { Story } = defineMeta({
		title: 'Board/BoardStats',
		component: BoardStats,
		args: { docCount: 74, codeCount: 53 },
	});
</script>

<Story name="Stats and legend" />

<Story name="Documentation only" args={{ codeCount: 0 }} />
