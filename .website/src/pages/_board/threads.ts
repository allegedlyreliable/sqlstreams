import type { CollectionEntry } from 'astro:content';
import { boards, stickyIds } from './navigation';

export type DocsEntry = CollectionEntry<'docs'>;

// one shape for every article a board can name; its URL is /<id>/
export type Thread = {
	id: string;
	group: string | null;
	title: string;
	description: string;
	// path relative to .website/, where the build runs
	filePath: string;
	entry: DocsEntry;
};

export function siteThreads(docs: DocsEntry[]): Thread[] {
	const threads = docs.map(docsThread);
	const ids = threads.map((thread) => thread.id);
	const visibleIds = new Set([...stickyIds, ...boards.flatMap((board) => board.threads(ids))]);
	return threads.filter((thread) => visibleIds.has(thread.id));
}

// the repo-rooted path GitHub links need
export function repositoryFilePath(thread: Thread): string {
	return `.website/${thread.filePath}`;
}

// ***************
// *** HELPERS ***
// ***************

function docsThread(entry: DocsEntry): Thread {
	return {
		id: entry.id,
		group: entry.data.group ?? null,
		title: entry.data.title,
		description: entry.data.description ?? entry.data.title,
		filePath: entryFilePath(entry),
		entry,
	};
}

function entryFilePath(entry: DocsEntry): string {
	if (entry.filePath === undefined) {
		throw new Error(`entry "${entry.id}" carries no filePath`);
	}
	return entry.filePath;
}
