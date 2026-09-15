<script lang="ts" module>
	import { defineMeta } from '@storybook/addon-svelte-csf';
	import GroupMap from './group-map.svelte';
	const { Story } = defineMeta({
		title: 'Board/GroupMap',
		component: GroupMap,
		args: {
			currentHref: '/guides/stop-a-consumer/',
			map: {
				title: 'Consumer',
				sections: [
					{
						title: 'Concepts',
						entries: [
							{
								title: 'Consumer groups',
								href: '/concepts/consumer-groups/',
							},
						],
					},
					{
						title: 'Guides',
						entries: [
							{
								title: 'Stop a consumer',
								href: '/guides/stop-a-consumer/',
							},
						],
					},
					{
						title: 'Reference',
						entries: [
							{
								title: 'Consume',
								href: '/reference/consume/',
							},
						],
					},
				],
			},
		},
	});
</script>

<Story name="CurrentGuide" />
<Story name="CurrentReference" args={{ currentHref: '/reference/consume/' }} />
