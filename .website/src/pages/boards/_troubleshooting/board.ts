import type { Board } from '../../_board/model';
import { isErrorThread } from '../../_board/identifiers';

export const board: Board = {
	title: 'Troubleshooting',
	slug: 'troubleshooting',
	description: 'errors, log events, metrics, and alerts — one thread per SQL code',
	threads: (ids) => ids.filter(isErrorThread).sort(),
};
