// @ts-check
import { URL } from 'node:url';
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';
import svelte from '@astrojs/svelte';
import expressiveCode from 'astro-expressive-code';
import { isSearchEngineIndexable } from './src/search-engine-index.ts';
import { siteUrl } from './src/site.ts';

// named keyword families only -- keyword.operator stays ink
const keywordScopes = [
	'keyword.control',
	'keyword.function',
	'keyword.package',
	'keyword.import',
	'keyword.type',
	'keyword.var',
	'keyword.const',
	'keyword.struct',
	'keyword.interface',
	'keyword.map',
	'keyword.channel',
	'keyword.other',
	'storage',
	'constant.language',
];

// https://astro.build/config
export default defineConfig({
	site: siteUrl,
	redirects: {
		'/errors/': '/boards/troubleshooting/',
		'/reference/': '/boards/reference/',
	},
	vite: {
		build: {
			// Vite 8's baseline-widely-available list, pinned -- the default
			// floats per Vite major, so the floor moves only by deliberate edit
			target: ['chrome111', 'edge111', 'firefox114', 'safari16.4', 'ios16.4'],
		},
		// PGlite locates its wasm assets itself; pre-bundling breaks the paths
		optimizeDeps: { exclude: ['@electric-sql/pglite'] },
	},
	integrations: [
		expressiveCode({
			themes: [
				{
					name: 'classic',
					type: 'light',
					// console-sql-pale, ink
					colors: {
						'editor.background': '#f6f9fb',
						'editor.foreground': '#22303c',
					},
					settings: [
						{ settings: { foreground: '#22303c' } },
						// band-blue-end
						{ scope: keywordScopes, settings: { foreground: '#184e7c', fontStyle: 'bold' } },
						// sticky-label-red
						{ scope: ['string'], settings: { foreground: '#b03a2e' } },
						// ink-faint
						{ scope: ['comment'], settings: { foreground: '#7c8b98', fontStyle: 'italic' } },
					],
				},
				{
					name: 'night',
					type: 'dark',
					// console-sql-pitch, ink-silver
					colors: {
						'editor.background': '#12161b',
						'editor.foreground': '#c8ccd2',
					},
					settings: [
						{ settings: { foreground: '#c8ccd2' } },
						// console-keyword-sky
						{ scope: keywordScopes, settings: { foreground: '#7fb3dd', fontStyle: 'bold' } },
						// sticky-label-coral
						{ scope: ['string'], settings: { foreground: '#e0685a' } },
						// ink-silver-faint
						{ scope: ['comment'], settings: { foreground: '#6b7178', fontStyle: 'italic' } },
					],
				},
			],
			themeCssSelector: (theme) => `[data-board-style='${theme.name}']`,
			useDarkModeMediaQuery: false,
			minSyntaxHighlightingColorContrast: 0,
			defaultProps: { frame: 'none' },
			frames: {
				extractFileNameFromCode: false,
				removeCommentsWhenCopyingTerminalFrames: false,
			},
			styleOverrides: {
				borderColor: 'var(--border-row)',
				borderRadius: '0',
				borderWidth: 'var(--space-1)',
				codeFontFamily: 'var(--font-database)',
				codeFontSize: 'var(--font-size-code)',
				codePaddingBlock: 'var(--space-13)',
				codePaddingInline: 'var(--space-16)',
				uiFontFamily: 'var(--font-chrome)',
				frames: { shadowColor: 'transparent' },
			},
		}),
		svelte(),
		mdx(),
		sitemap({
			filter: (/** @type {string} */ page) => isSearchEngineIndexable(new URL(page).pathname),
		}),
	],
});
