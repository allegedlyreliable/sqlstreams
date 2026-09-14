// how deeply the author reviewed an LLM-drafted thread; none is a thread
// the author checked by hand and renders no notice
export const slopLevels = ['none', 'low', 'medium', 'high'] as const;

export type SlopLevel = (typeof slopLevels)[number];

export type SlopNoticeLevel = Exclude<SlopLevel, 'none'>;

export type SlopCopy = {
	// filled blocks on the three-block meter
	blocks: number;
	body: string;
};

export const slopCopy: Record<SlopNoticeLevel, SlopCopy> = {
	low: {
		blocks: 1,
		body: 'An LLHuMan wrote this. I have thorougly reviewed this. The good parts are mine. The bad parts are not mine.',
	},
	medium: {
		blocks: 2,
		body: 'An LLM wrote this. I have briefly reviewed this. Content is accurate details might not be.',
	},
	high: {
		blocks: 3,
		body: "An LLM wrote this. I have not reviewed it yet. Content should be accurate enough but you've been warned",
	},
};
