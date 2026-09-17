// one page-colored square laid over the text, gone again after a short life
export type SpawnedSquare = {
	id: number;
	// the top-left corner, viewport px
	x: number;
	y: number;
	// the side length, px
	size: number;
};

// one strip of repeated personal text crossing the viewport along its own
// direction: it rolls in from one edge, then loops forever
export type SpawnedLine = {
	id: number;
	// the viewport point the strip passes through, px
	x: number;
	y: number;
	// the strip's direction, degrees clockwise from rightward
	angle: number;
	fontSize: number;
	// px per second along the strip
	speed: number;
	// from the point back past the edge the text rolls in from, px
	entry: number;
	// the roll-in travel, edge to edge, px
	distance: number;
	copy: string;
	// one copy's rendered width, px; 0 until measured
	copyWidth: number;
};
