import type { Board } from '../../_board/model';
import { isDecisionRecordThread } from '../../_board/identifiers';

export const board: Board = {
	title: 'Decision records',
	slug: 'decisions',
	description: 'the why behind shipped behavior — every settled design decision, append-only',
	// newest decision records first
	threads: (ids) => ids.filter(isDecisionRecordThread).sort().reverse(),
};
