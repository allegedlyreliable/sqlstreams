<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import ResourceCard from './resource-card.svelte';
	const { Story } = defineMeta({
		title: 'Board/Resource Card',
		component: ResourceCard,
		args: { data: { title: 'Group: charge-cards', lines: ['Completed: 1 · 2 · 3 · 4'] } },
	});
</script>

<Story name="One line">
	{#snippet template(args)}<div style="width: 250px; height: 88px">
			<ResourceCard {...args} />
		</div>{/snippet}
</Story>
<Story
	name="Two lines"
	args={{
		data: { title: 'Group: audit-payments', lines: ['Completed: 1 · 2', 'Remaining: 3 · 4'] },
	}}
>
	{#snippet template(args)}<div style="width: 250px; height: 88px">
			<ResourceCard {...args} />
		</div>{/snippet}
</Story>
<Story name="Container" args={{ data: { title: 'Group: charge-cards', lines: null } }}>
	{#snippet template(args)}<div style="width: 250px; height: 160px">
			<ResourceCard {...args} />
		</div>{/snippet}
</Story>
