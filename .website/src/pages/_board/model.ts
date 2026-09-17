export type Board = {
	grouped: boolean;
	// false drops the board from the header nav under the 640px layout collapse
	inPhoneNav: boolean;
	title: string;
	slug: string;
	description: string;
	threads: (ids: string[]) => string[];
};

export type BoardRowData = {
	title: string;
	href: string;
	description: string;
	threadCount: number;
	lastPostTitle: string;
	lastPostHref: string;
	lastPostDate: string;
	scopeHrefs: string[];
};

export type ThreadRowData = {
	metadata: string | null;
	title: string;
	href: string;
	lastUpdatedDate: string;
};

export type StickyRowData = {
	title: string;
	href: string;
	lastUpdatedDate: string;
};
