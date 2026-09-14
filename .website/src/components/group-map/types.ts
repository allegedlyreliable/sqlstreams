export type GroupMapEntry = {
	title: string;
	href: string;
};

export type GroupMapSection = {
	title: string;
	entries: GroupMapEntry[];
};

export type GroupMapData = {
	title: string;
	sections: GroupMapSection[];
};
