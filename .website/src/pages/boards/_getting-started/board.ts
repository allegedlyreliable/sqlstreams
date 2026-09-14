import type { Board } from '../../_board/model';

export const board: Board = {
	title: 'Getting Started',
	slug: 'getting-started',
	description: 'first produce and consume, why SQLStreams, benchmarks, and the roadmap',
	threads: () => ['quickstart', 'why-sqlstreams', 'benchmarks', 'roadmap'],
};
