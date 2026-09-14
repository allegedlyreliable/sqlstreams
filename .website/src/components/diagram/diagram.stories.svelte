<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import Diagram from './diagram.svelte';
	import { layoutResources } from './layout';

	const independent = await layoutResources(
		[
			{
				id: 'stream',
				title: 'Stream: payments.requested',
				width: 250,
				lines: ['Messages: 1 · 2 · 3 · 4'],
			},
			{
				id: 'charges',
				title: 'Group: charge-cards',
				width: 190,
				lines: ['Completed: 1 · 2 · 3 · 4'],
			},
			{
				id: 'audit',
				title: 'Group: audit-payments',
				width: 190,
				lines: ['Completed: 1 · 2', 'Remaining: 3 · 4'],
			},
		],
		[
			['stream', 'charges'],
			['stream', 'audit'],
		],
	);
	const nested = await layoutResources(
		[
			{
				id: 'group',
				title: 'Group: charge-cards',
				width: 0,
				children: [
					{ id: 'A', title: 'Instance A', width: 140, lines: ['Claims: 1–2'] },
					{ id: 'B', title: 'Instance B', width: 140, lines: ['Claims: 3–4'] },
				],
			},
		],
		[],
	);
	const { Story } = defineMeta({
		title: 'Board/Diagram',
		component: Diagram,
		args: {
			graph: independent,
			description: 'Two groups keep independent progress.',
			caption: 'The messages are stored once. Each group tracks its own progress.',
		},
	});
</script>

<Story name="Independent progress" />
<Story
	name="Nested instances"
	args={{
		graph: nested,
		description: 'Two consumer instances inside one group.',
		caption: 'Instances share the group’s work.',
	}}
/>
