import ELK from 'elkjs/lib/elk.bundled.js';
import type { ElkNode } from 'elkjs';
import { Position } from '@xyflow/svelte';
import type { Edge } from '@xyflow/svelte';
import { cardGeometry } from '../resource-card/geometry';
import type { DiagramLayout, ResourceNode } from './types';

type LayoutNode = ElkNode & { node: ResourceNode; children: LayoutNode[] };

const elk = new ELK();
const graphLayoutOptions = {
	'elk.algorithm': 'layered',
	'elk.direction': 'DOWN',
	'elk.layered.nodePlacement.bk.fixedAlignment': 'BALANCED',
	'elk.padding': '[top=16,left=16,bottom=16,right=16]',
	'elk.spacing.nodeNode': '18',
	'elk.layered.spacing.nodeNodeBetweenLayers': '48',
};

function layoutNodes(nodes: ResourceNode[], parentId: string | undefined): LayoutNode[] {
	return nodes
		.filter((node) => node.parentId === parentId)
		.map((node) => {
			const layout: LayoutNode = { id: node.id, node, children: layoutNodes(nodes, node.id) };
			if (node.data.lines !== null) {
				if (node.width === undefined || node.height === undefined) {
					throw new Error(`Diagram resource "${node.id}" has no dimensions`);
				}
				layout.width = node.width;
				layout.height = node.height;
				return layout;
			}
			const padding = cardGeometry.childrenPadding;
			layout.layoutOptions = {
				'elk.algorithm': 'box',
				'elk.aspectRatio': '4',
				'elk.padding': `[top=${cardGeometry.headerHeight + padding},left=${padding},bottom=${padding},right=${padding}]`,
			};
			return layout;
		});
}

function positionedNodes(layout: LayoutNode): ResourceNode[] {
	const { x, y, width, height } = layout;
	if (x === undefined || y === undefined || width === undefined || height === undefined) {
		throw new Error(`Diagram resource "${layout.id}" has no complete layout`);
	}
	return [
		{
			...layout.node,
			position: { x, y },
			width,
			height,
			handles: [
				{ type: 'target', position: Position.Top, x: width / 2, y: 0 },
				{ type: 'source', position: Position.Bottom, x: width / 2, y: height },
			],
		},
		...layout.children.flatMap(positionedNodes),
	];
}

export async function layoutNodesAndEdges(
	nodes: ResourceNode[],
	edges: Edge[],
): Promise<DiagramLayout> {
	const input: ElkNode & { children: LayoutNode[] } = {
		id: 'diagram',
		layoutOptions: graphLayoutOptions,
		children: layoutNodes(nodes, undefined),
		edges: edges.map((edge) => ({ id: edge.id, sources: [edge.source], targets: [edge.target] })),
	};
	const result = await elk.layout(input);
	if (!result.children || result.width === undefined || result.height === undefined) {
		throw new Error('Diagram has no complete layout');
	}
	return {
		nodes: result.children.flatMap(positionedNodes),
		edges,
		width: result.width,
		height: result.height,
	};
}
