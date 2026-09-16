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
			{ id: 'consumer-groups', title: 'Consumer groups', rows },
			{
				id: 'messages',
				title: 'Messages',
				rows: [
					{
						metadata: null,
						lastUpdatedDate: '2026-09-14',
						title: 'Message metadata',
						href: '/reference/message-metadata/',
					},
				],
			},
		],
	}}
/>
