import type { Board } from '../../_board/model';

export const board: Board = {
	title: 'Concepts',
	slug: 'concepts',
	description:
		'queue & log, the shape of the API, lifecycle, handler outcomes, consumer group config, message key, ordering, routing, fan-out, architecture, table design',
	threads: () => [
		'concepts/queue-and-log',
		'concepts/api-shape',
		'concepts/fan-out',
		'concepts/lifecycle',
		'concepts/handler-outcomes',
		'concepts/consumer-group-config',
		'concepts/routing',
		'concepts/message-key',
		'concepts/ordering',
		'concepts/architecture',
		'concepts/table-design',
	],
};
