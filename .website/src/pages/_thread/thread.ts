import { groupTitle } from '../_board/groups';
import { lastCommitDate } from '../../helpers/last-commit-date';
import { repositoryUrl } from '../../site';
import type { ThreadLink } from '../../components/prev-next/types';
import type { Board } from '../_board/model';
import { boards, boardHref, stickyIds } from '../_board/navigation';
import { repositoryFilePath, type Thread } from '../_board/threads';
import type { ThreadData } from './model';
import { groupMap } from './group-map';

export function threadData(thread: Thread, threads: Thread[]): ThreadData {
	const ids = threads.map((candidate) => candidate.id);
	const board = boards.find((candidate) => candidate.threads(ids).includes(thread.id)) ?? null;
	if (board === null && !stickyIds.includes(thread.id)) {
		throw new Error(`thread "${thread.id}" belongs to neither a board nor the homepage stickies`);
	}

	const members = board === null ? stickyIds : board.threads(ids);
	const position = members.indexOf(thread.id);
	const previousThread = threads.find((candidate) => candidate.id === members[position - 1]);
	const nextThread = threads.find((candidate) => candidate.id === members[position + 1]);

	return {
		groupMap: thread.group === null ? null : groupMap(thread.group, threads),
		group: toGroupLink(board, thread.group),
		board: board === null ? null : { title: board.title, href: boardHref(board) },
		postedDate: lastCommitDate(thread.filePath),
		postCount: threads.length,
		editHref: `${repositoryUrl}/edit/main/${repositoryFilePath(thread)}`,
		reportHref: `${repositoryUrl}/issues/new?title=${encodeURIComponent(`docs: ${thread.title}`)}`,
		previous: toThreadLink(previousThread),
		next: toThreadLink(nextThread),
		slop: thread.entry.data.slop ?? null,
	};
}

// ***************
// *** HELPERS ***
// ***************

function toThreadLink(thread: Thread | undefined): ThreadLink | null {
	if (thread === undefined) {
		return null;
	}
	return { title: thread.title, href: `/${thread.id}/` };
}

function toGroupLink(board: Board | null, group: string | null): ThreadLink | null {
	if (board === null || !board.grouped || group === null) {
		return null;
	}
	return { title: groupTitle(group), href: `${boardHref(board)}#${group}` };
}
