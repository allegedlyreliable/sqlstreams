import { isErrorThread } from './identifiers';
import type { SiteStats } from './model';
import type { Thread } from './threads';

export function siteStats(threads: Thread[]): SiteStats {
	return {
		docCount: threads.length,
		codeCount: threads.filter((thread) => isErrorThread(thread.id)).length,
	};
}
