import { threadRows } from '../../_board/rows';
import type { Thread } from '../../_board/threads';
import { board } from './board';
import type { ReferenceSectionRows } from './model';
import { sections } from './sections';

export function referenceRows(threads: Thread[]): ReferenceSectionRows[] {
	return sections.map((section) => ({
		id: section.id,
		title: section.title,
		rows: threadRows({ ...board, threads: () => section.ids }, threads),
	}));
}
