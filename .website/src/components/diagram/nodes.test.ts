import { describe, expect, it } from 'vitest';
import { resourceNodes } from './nodes';

describe('resourceNodes', () => {
	it('grows sibling cards together when one body needs more room', () => {
		const nodes = resourceNodes(
			[
				{ id: 'charges', title: 'Group: charge-cards', width: 190, lines: ['Completed: 1–4'] },
				{
					id: 'audit',
					title: 'Group: audit-payments',
					width: 190,
					lines: ['Completed: 1–2', 'Remaining: 3–4', 'Instance A stopped', 'Instance B running'],
				},
			],
			null,
		);
		const shortCard = resourceNodes(
			[{ id: 'charges', title: 'Group: charge-cards', width: 190, lines: ['Completed: 1–4'] }],
			null,
		)[0]!;
		expect(nodes[0]!.height).toBe(nodes[1]!.height);
		expect(nodes[0]!.height).toBeGreaterThan(shortCard.height!);
	});

	it('places parents before their instances and retains each resource’s content', () => {
		const nodes = resourceNodes(
			[
				{
					id: 'charges',
					title: 'Group: charge-cards',
					width: 0,
					children: [{ id: 'A', title: 'Instance A', width: 140, lines: ['Claims: 1–2'] }],
				},
			],
			null,
		);
		expect(nodes.map((node) => node.id)).toEqual(['charges', 'A']);
		expect(nodes[0]!.parentId).toBeUndefined();
		expect(nodes[1]!.parentId).toBe('charges');
		expect(nodes[1]!.data).toEqual({ title: 'Instance A', lines: ['Claims: 1–2'] });
	});
});
