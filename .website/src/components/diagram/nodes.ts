import { sharedCardHeight } from '../resource-card/geometry';
import type { DiagramResource, ResourceNode } from './types';

export function resourceNodes(
	resources: DiagramResource[],
	parentId: string | null,
): ResourceNode[] {
	const height = sharedCardHeight(
		resources.map((resource) => ('lines' in resource ? resource.lines.length : 0)),
	);
	return resources.flatMap((resource) => {
		const node: ResourceNode = {
			id: resource.id,
			type: 'resource',
			position: { x: 0, y: 0 },
			data: { title: resource.title, lines: 'lines' in resource ? resource.lines : null },
		};
		if (parentId !== null) node.parentId = parentId;
		if ('children' in resource) return [node, ...resourceNodes(resource.children, resource.id)];
		return [{ ...node, width: resource.width, height }];
	});
}
