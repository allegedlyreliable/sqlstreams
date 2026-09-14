<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import BoardGroup from './board-group.svelte';
	const rows = [
		{
			title: 'Consume',
			href: '/reference/consume/',
			metadata: null,
			lastUpdatedDate: '2026-09-14',
		},
	];
	const { Story } = defineMeta({
		title: 'Board/BoardGroup',
		component: BoardGroup,
		args: { id: 'consumer', title: 'Consumer', sections: [{ id: null, title: null, rows }] },
	});
</script>

<Story name="Articles" />
<Story
	name="Subsections"
	args={{
		sections: [
			{ id: 'go-api', title: 'Go API', rows },
			{
				id: 'configuration',
				title: 'Configuration',
				rows: [
					{
						metadata: null,
						lastUpdatedDate: '2026-09-14',
						title: 'Consumer group settings',
						href: '/reference/consumer-group-settings/',
					},
				],
			},
		],
	}}
/>
