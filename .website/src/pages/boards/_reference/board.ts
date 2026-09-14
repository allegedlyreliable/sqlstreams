import type { Board } from '../../_board/model';

export const board: Board = {
	title: 'Reference',
	slug: 'reference',
	description:
		'one thread per handle and instance — every verb and every config field with its default, checked against the shipped library',
	// the handles in the order a program meets them
	threads: () => [
		'reference/client',
		'reference/pool',
		'reference/stream',
		'reference/maintenance',
		'reference/producer',
		'reference/consumer',
		'reference/key',
		'reference/scheduler',
		'reference/system',
		'reference/manager',
		'reference/metrics',
		'reference/alerts',
		'reference/message-options',
		'reference/diagnostics',
	],
};
