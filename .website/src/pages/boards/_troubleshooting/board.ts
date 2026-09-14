import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: false,
	title: 'Troubleshooting',
	slug: 'troubleshooting',
	description: 'Start with a symptom or code, identify its cause, and recover.',
	threads: () => ['troubleshooting/queued-message-cannot-start'],
};
