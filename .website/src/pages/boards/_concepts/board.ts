import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: true,
	title: 'Concepts',
	slug: 'concepts',
	description: 'Understand how delivery, shared groups, and leases behave.',
	threads: () => ['concepts/consumer-groups', 'concepts/consumer-leases'],
};
