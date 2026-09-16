<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import PrevNext from './prev-next.svelte';

	const { Story } = defineMeta({
		title: 'Board/PrevNext',
		component: PrevNext,
		args: {
			previous: {
				title: 'Retries and failed messages',
				href: '/concepts/retries-and-failed-messages/',
			},
			next: { title: 'Ordering and concurrency', href: '/concepts/ordering/' },
		},
	});
</script>

<Story name="Both" />

<Story name="LastThread" args={{ next: null }} />

<Story name="FirstThread" args={{ previous: null }} />
