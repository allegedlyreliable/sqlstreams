<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import GroupMapSection from './group-map-section.svelte';
	const { Story } = defineMeta({
		title: 'Board/GroupMapSection',
		component: GroupMapSection,
		args: {
			currentHref: null,
			section: {
				title: 'Guides',
				entries: [
					{
						title: 'Stop a consumer',
						href: '/guides/stop-a-consumer/',
					},
					{
						title: 'Size a consumer queue',
						href: '/guides/size-a-consumer-queue/',
					},
				],
			},
		},
	});
</script>

<Story name="Collapsed" />
<Story name="CurrentPage" args={{ currentHref: '/guides/stop-a-consumer/' }} />
