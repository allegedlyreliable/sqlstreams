import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: true,
	title: 'Guides',
	slug: 'guides',
	description: 'Complete a task in an existing application and check the result.',
	threads: () => ['guides/stop-a-consumer', 'guides/size-a-consumer-queue'],
};
