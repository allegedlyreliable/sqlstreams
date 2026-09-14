import type { GroupMapData } from '../../components/group-map/types';
import type { Thread } from '../_board/threads';
import { boards, stickyIds } from '../_board/navigation';
import { groupTitle } from '../_board/groups';

export function groupMap(group: string, threads: Thread[]): GroupMapData {
	const ids = threads.map((thread) => thread.id);
	const members = threads.filter((thread) => thread.group === group);
	const groups = [
		{ title: 'Start here', ids: stickyIds },
		...boards.map((board) => ({ title: board.title, ids: board.threads(ids) })),
	];
	return {
		title: groupTitle(group),
		sections: groups
			.map((group) => ({
				title: group.title,
				entries: group.ids.flatMap((id) => {
					const thread = members.find((candidate) => candidate.id === id);
					return thread === undefined ? [] : [{ title: thread.title, href: `/${thread.id}/` }];
				}),
			}))
			.filter((section) => section.entries.length > 0),
	};
}
