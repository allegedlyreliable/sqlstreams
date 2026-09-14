<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import ThreadRow from './thread-row.svelte';

	const { Story } = defineMeta({
		title: 'Board/ThreadRow',
		component: ThreadRow,
		args: {
			index: 0,
			decisionDate: null,
			metadata: null,
			title: 'Ordering & Concurrency',
			href: '/concepts/ordering/',
			updated: false,
			lastUpdatedDate: '2026-08-23',
			onVisit: () => {},
		},
	});
</script>

<Story name="Read" />
<Story name="Updated" args={{ updated: true }} />
<Story
	name="Error thread"
	args={{
		index: 1,
		title: 'stream not found [SQL0005]',
		href: '/errors/SQL0005/',
		metadata: 'permanent',
	}}
/>

<Story
	name="Decision record"
	args={{
		title: '0798 — Threads declare a slop level',
		href: '/decisions/0798/',
		decisionDate: '2026-09-13',
		metadata: 'accepted',
	}}
/>
<Story name="Superseded decision" args={{ decisionDate: '2026-09-12', metadata: 'superseded' }} />
<Story name="Rejected decision" args={{ decisionDate: '2026-09-12', metadata: 'rejected' }} />
