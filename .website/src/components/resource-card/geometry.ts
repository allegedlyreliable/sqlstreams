export const cardGeometry = {
	headerHeight: 28,
	lineHeight: 18,
	bodyPadding: 8,
	borderWidth: 1,
	minimumHeight: 88,
	childrenPadding: 16,
};

export function sharedCardHeight(lineCounts: number[]): number {
	const lines = Math.max(0, ...lineCounts);
	const bodyHeight = lines * cardGeometry.lineHeight + 2 * cardGeometry.bodyPadding;
	return Math.max(
		cardGeometry.minimumHeight,
		cardGeometry.headerHeight + bodyHeight + 2 * cardGeometry.borderWidth,
	);
}
