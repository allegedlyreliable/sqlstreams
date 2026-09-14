import { codeRecords } from '../../data/codes';
import { isErrorThread, threadCode } from './identifiers';
import type { DocsEntry } from './threads';

export function codeMetadata(entry: DocsEntry): string | null {
	if (!isErrorThread(entry.id)) return null;
	const record = codeRecords[threadCode(entry.id)];
	if (record === undefined) throw new Error(`code page "${entry.id}" has no declaration`);
	switch (record.kind) {
		case 'error':
			return record.recovery;
		case 'event': {
			if (entry.data.level === undefined) throw new Error(`event page "${entry.id}" has no level`);
			return entry.data.level;
		}
		case 'metric':
			return record.metric_kind;
		case 'alert':
			return record.severity;
	}
}
