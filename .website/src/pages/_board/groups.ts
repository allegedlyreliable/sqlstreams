import type { BoardGroupSection } from '../../components/board-group/types';
import type { Thread } from './threads';

export type BoardGroupData = {
	id: string;
	title: string;
	sections: BoardGroupSection[];
};

export function groupTitle(groupId: string): string {
	const words = groupId.replaceAll('-', ' ');
	return words.charAt(0).toUpperCase() + words.slice(1);
}

export function groupRows(sections: BoardGroupSection[], threads: Thread[]): BoardGroupData[] {
	const groupIdsByHref = new Map(threads.map((thread) => [`/${thread.id}/`, thread.group]));
	const groups = new Map<string, BoardGroupData>();
	for (const section of sections) {
		for (const row of section.rows) {
			const groupId = groupIdsByHref.get(row.href);
			if (groupId === undefined || groupId === null) {
				throw new Error(`Article "${row.href}" needs group frontmatter to appear in a board group`);
			}
			let group = groups.get(groupId);
			if (group === undefined) {
				group = { id: groupId, title: groupTitle(groupId), sections: [] };
				groups.set(groupId, group);
			}
			let groupSection = group.sections.find((candidate) => candidate.id === section.id);
			if (groupSection === undefined) {
				groupSection = { id: section.id, title: section.title, rows: [] };
				group.sections.push(groupSection);
			}
			groupSection.rows.push(row);
		}
	}
	return [...groups.values()];
}
