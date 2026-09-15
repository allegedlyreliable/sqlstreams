import type { Edge, Node } from '@xyflow/svelte';
import type { QueueNode } from '../queue-card/types';

export type ResourceNode = Node<{ title: string; lines: string[] | null }, 'resource'>;

export type DiagramLayout = {
	nodes: (ResourceNode | QueueNode)[];
	edges: Edge[];
	width: number;
	height: number;
};

export type DiagramResource = {
	id: string;
	title: string;
	width: number;
} & ({ lines: string[] } | { children: DiagramResource[] });

export type DiagramConnection = [from: string, to: string];
