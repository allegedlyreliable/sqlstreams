import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: true,
	inPhoneNav: true,
	title: 'Concepts',
	slug: 'concepts',
	description:
		'Understand how SQLStreams stores messages and coordinates processing and maintenance.',
	threads: () => [
		'concepts/streams-and-messages',
		'concepts/consumer-groups',
		'concepts/consumer-leases',
		'concepts/retries-and-failed-messages',
		'concepts/ordering',
		'concepts/background-managers-and-workers',
	],
};
