import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';
import { slopLevels } from './components/slop-notice/types';

export const collections = {
	docs: defineCollection({
		loader: glob({
			pattern: '**/*.{md,mdx}',
			base: './src/content/docs',
			// ids and their URLs keep the source file path's casing
			generateId: ({ entry }) => entry.replace(/\.(md|mdx)$/, ''),
		}),
		schema: z.object({
			group: z
				.string()
				.regex(/^[a-z]+(?:-[a-z]+)*$/)
				.optional()
				.describe('the documentation group used by boards, related pages, and breadcrumbs'),
			title: z.string().describe('the thread title; the H1 the page renders under'),
			description: z.string().optional().describe('the meta description for the article'),
			slop: z
				.enum(slopLevels)
				.optional()
				.describe(
					'how far the author reviewed the LLM-drafted body; none is a hand-checked thread and renders no notice; absent only on the diagnostics reference',
				),
		}),
	}),
};
