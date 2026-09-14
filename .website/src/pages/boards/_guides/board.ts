import type { Board } from '../../_board/model';

export const board: Board = {
	title: 'Guides',
	slug: 'guides',
	description:
		'transactional produce, side effects & retries, replay, dead letters, where a new group starts, ordered delivery, consumer tuning, consumer timeouts, schema versions, migrations, schedules',
	threads: () => [
		'guides/transactional-produce',
		'guides/side-effects-and-retries',
		'guides/replay',
		'guides/dead-letters',
		'guides/new-group-start',
		'guides/ordered-delivery',
		'guides/consumer-tuning',
		'guides/consumer-timeouts',
		'guides/schema-versions',
		'guides/migrations',
		'guides/schedules',
	],
};
