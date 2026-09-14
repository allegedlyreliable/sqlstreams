import type { ThreadRowData } from '../../_board/model';

export type ReferenceSection = {
	id: string;
	title: string;
	ids: string[];
};

export type ReferenceSectionRows = {
	id: string;
	title: string;
	rows: ThreadRowData[];
};
