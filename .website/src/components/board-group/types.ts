import type { ThreadRowData } from '../../pages/_board/model';
export type BoardGroupSection = {
	id: string | null;
	title: string | null;
	rows: ThreadRowData[];
};
