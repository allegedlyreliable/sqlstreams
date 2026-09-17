import type { Board } from '../../_board/model';
import { sections } from './sections';

export const board: Board = {
	grouped: true,
	inPhoneNav: false,
	title: 'Reference',
	slug: 'reference',
	description: 'Go API, configuration, CLI, metrics, logs, and alerts — look up an exact contract',
	threads: () => sections.flatMap((section) => section.ids),
};
