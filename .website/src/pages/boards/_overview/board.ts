import type { Board } from '../../_board/model';

export const board: Board = {
	grouped: false,
	title: 'Overview',
	slug: 'overview',
	description: 'Start here: what SQLStreams is, how fast it runs, and where it is going.',
	threads: () => ['quickstart', 'why-sqlstreams', 'benchmarks', 'roadmap'],
};
