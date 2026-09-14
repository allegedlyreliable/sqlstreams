export type Board = {
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
	decisionDate: string | null;
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

export type SiteStats = {
	docCount: number;
	codeCount: number;
	decisionRecordCount: number;
};
