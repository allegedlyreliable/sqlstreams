import type { BoardSectionData } from './model';
import { threadRows } from '../../_board/rows';
import type { Thread } from '../../_board/threads';
import { board } from './board';

export function troubleshootingSections(threads: Thread[]): BoardSectionData[] {
	return [
		{ kind: 'error', title: 'Errors', label: 'Recovery' },
		{ kind: 'event', title: 'Log events', label: 'Level' },
		{ kind: 'metric', title: 'Metrics', label: 'Kind' },
		{ kind: 'alert', title: 'Alerts', label: 'Severity' },
	].map(({ kind, title, label }) => ({
		title,
		columnLabels: { threadCount: label, lastPost: 'Updated' },
		rows: threadRows(
			board,
			threads.filter(
				(thread) => thread.entry.collection === 'docs' && thread.entry.data.kind === kind,
			),
		),
	}));
}
