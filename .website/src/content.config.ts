import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';
import { slopLevels } from './components/slop-notice/types';

export const collections = {
	docs: defineCollection({
		loader: glob({
			pattern: '**/*.{md,mdx}',
			base: './src/content/docs',
			// the default generateId slugifies (SQL0005 -> sql0005); ids and the
			// URLs built from them keep the file path's own casing
			generateId: ({ entry }) => entry.replace(/\.(md|mdx)$/, ''),
		}),
		schema: z.object({
			group: z
				.string()
				.regex(/^[a-z]+(?:-[a-z]+)*$/)
				.optional()
				.describe('the documentation group used by boards, related pages, and breadcrumbs'),
			title: z.string().describe('the thread title; the H1 the page renders under'),
			description: z
				.string()
				.optional()
				.describe('the meta description for ordinary threads; code pages derive theirs'),
			kind: z
				.enum(['error', 'event', 'metric', 'alert'])
				.optional()
				.describe('what the SQL code names; present on every error-board page, absent elsewhere'),
			recovery: z
				.enum(['permanent', 'transient'])
				.optional()
				.describe('whether an unchanged retry can succeed; error pages only'),
			level: z
				.enum(['info', 'warn', 'error'])
				.optional()
				.describe('the log level the event is emitted at; event pages only'),
			consequence: z
				.string()
				.optional()
				.describe('the consequence clause shown beside the code, verbatim from the declaration'),
			fix: z
				.string()
				.optional()
				.describe('the declared fix, verbatim; absent when the code cannot know one'),
			slop: z
				.enum(slopLevels)
				.optional()
				.describe(
					'how far the author reviewed the LLM-drafted body; none is a hand-checked thread and renders no notice; absent only on the diagnostics reference and the code threads',
				),
		}),
	}),
};
