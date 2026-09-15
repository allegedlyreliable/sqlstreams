import type { DiagramLayout } from '../diagram/types';
import type { QueueSnapshot } from './types';

export function queueDiagram(snapshot: QueueSnapshot): DiagramLayout {
	return {
		width: 526,
		height: 470,
		nodes: [
			{
				id: 'queue',
				type: 'queue',
				data: snapshot,
				position: { x: 16, y: 16 },
				width: 494,
				height: 438,
			},
		],
		edges: [],
	};
}
