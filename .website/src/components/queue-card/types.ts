import type { Node } from '@xyflow/svelte';

export type QueueSnapshot = {
	instance: string;
	group: string;
	stream: string;
	processing: number;
	waiting: number[];
};

export type QueueNode = Node<QueueSnapshot, 'queue'>;
