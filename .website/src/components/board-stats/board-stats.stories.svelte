<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import BoardStats from './board-stats.svelte';

	const { Story } = defineMeta({
		title: 'Board/BoardStats',
		component: BoardStats,
		args: { docCount: 94 },
	});
</script>

<Story name="Stats and legend" />
