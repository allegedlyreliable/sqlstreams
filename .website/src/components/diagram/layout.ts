import { layoutNodesAndEdges } from './elk-layout';
import { resourceNodes } from './nodes';
import type { DiagramConnection, DiagramLayout, DiagramResource } from './types';

export function layoutResources(
	resources: DiagramResource[],
	connections: DiagramConnection[],
): Promise<DiagramLayout> {
	const nodes = resourceNodes(resources, null);
	const edges = connections.map(([source, target], index) => ({
		id: `connection-${index}`,
		source,
		target,
	}));
	return layoutNodesAndEdges(nodes, edges);
}
