import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: false,
	title: 'Troubleshooting',
	slug: 'troubleshooting',
	description: 'Start with a symptom or code, identify its cause, and recover.',
	threads: () => [
		'troubleshooting/queued-message-cannot-start',
		'troubleshooting/messages-not-claimed',
		'troubleshooting/messages-dead-lettered',
		'troubleshooting/consumer-handler-exceeds-timeout',
		'troubleshooting/lease-repeatedly-reclaimed',
		'troubleshooting/commit-confirmation-lost',
		'troubleshooting/produce-is-slow',
		'troubleshooting/schema-version-mismatch',
		'troubleshooting/migration-lock-timeout',
		'troubleshooting/worker-not-running',
	],
};
