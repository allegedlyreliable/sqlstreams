import type { Board } from '../../_board/model';

export const board: Board = {
	title: 'Compare',
	slug: 'compare',
	description: 'Kafka · RabbitMQ & SQS — shipped behavior only, no wishful checkmarks',
	threads: () => ['compare/kafka', 'compare/rabbitmq-sqs', 'compare/job-queues'],
};
