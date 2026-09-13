<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import SlopNotice from './slop-notice.svelte';

	const { Story } = defineMeta({
		title: 'Board/SlopNotice',
		component: SlopNotice,
		args: { level: 'high' },
	});
</script>

<Story name="Low" args={{ level: 'low' }}>
	{#snippet template(args)}
		<div class="thread-body">
			<SlopNotice {...args} />
		</div>
	{/snippet}
</Story>

<Story name="Medium" args={{ level: 'medium' }}>
	{#snippet template(args)}
		<div class="thread-body">
			<SlopNotice {...args} />
		</div>
	{/snippet}
</Story>

<Story name="High">
	{#snippet template(args)}
		<div class="thread-body">
			<SlopNotice {...args} />
		</div>
	{/snippet}
</Story>
