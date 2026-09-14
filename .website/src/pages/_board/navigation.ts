import type { JumpTarget } from '../../components/jump-to/types';
import type { Board } from './model';
import { board as gettingStarted } from '../boards/_getting-started/board';
import { board as concepts } from '../boards/_concepts/board';
import { board as guides } from '../boards/_guides/board';
import { board as reference } from '../boards/_reference/board';
import { board as troubleshooting } from '../boards/_troubleshooting/board';
import { board as compare } from '../boards/_compare/board';
import { board as decisions } from '../boards/_decisions/board';

// Shared navigation and membership; each board page owns its presentation.
export const boards: Board[] = [
	gettingStarted,
	concepts,
	guides,
	reference,
	troubleshooting,
	compare,
	decisions,
];

export const stickyIds = ['quickstart', 'why-sqlstreams'];

// the Jump to select navigates to each board's listing page
export const jumpTargets: JumpTarget[] = boards.map((board) => ({
	label: board.title,
	href: boardHref(board),
}));

export function boardHref(board: Board): string {
	return `/boards/${board.slug}/`;
}
