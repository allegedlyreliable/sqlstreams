import type { GroupMapData } from '../../components/group-map/types';
import type { ThreadLink } from '../../components/prev-next/types';
import type { SlopLevel } from '../../components/slop-notice/types';

export type ThreadData = {
	board: ThreadLink | null;
	group: ThreadLink | null;
	groupMap: GroupMapData | null;
	postedDate: string;
	postCount: number;
	editHref: string;
	reportHref: string;
	previous: ThreadLink | null;
	next: ThreadLink | null;
	slop: SlopLevel | null;
};
