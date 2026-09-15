<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import QueueCard from './queue-card.svelte';
	const { Story } = defineMeta({
		title: 'Board/Queue Card',
		component: QueueCard,
		args: {
			data: {
				instance: 'A',
				group: 'charge-cards',
				stream: 'payments.requested',
				processing: 1,
				waiting: [2, 3, 4],
			},
		},
	});
</script>

<Story name="Processing and waiting">
	{#snippet template(args)}
		<div style="width: 494px; height: 438px"><QueueCard {...args} /></div>
	{/snippet}
</Story>
