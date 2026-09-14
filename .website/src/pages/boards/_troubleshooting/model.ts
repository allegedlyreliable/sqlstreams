import type { ThreadRowData } from '../../_board/model';

export type BoardSectionData = {
	title: string;
	columnLabels: { threadCount: string | null; lastPost: string };
	rows: ThreadRowData[];
};
