import { describe, expect, it } from 'vitest';
import config from '../../astro.config.mjs';
import { codeRecords } from './codes';

describe('diagnostic declarations and documentation', () => {
	it('declares every code as one of the four kinds', () => {
		const unknown = Object.values(codeRecords)
			.filter((record) => !['error', 'event', 'metric', 'alert'].includes(record.kind))
			.map((record) => `${record.code} is kind ${record.kind}`);
		expect(unknown).toEqual([]);
	});

	it('redirects every diagnostic URL to Reference or Troubleshooting', () => {
		const redirects = config.redirects ?? {};
		const codePaths = Object.keys(codeRecords).map((code) => `/errors/${code}/`);
		const redirectedCodes = Object.keys(redirects).filter((path) =>
			/^\/errors\/SQL\d+\/$/.test(path),
		);
		expect(redirectedCodes.sort()).toEqual(codePaths.sort());
		for (const path of codePaths) {
			expect(redirects[path], path).toMatch(/^\/(reference|troubleshooting)\//);
		}
	});
});
