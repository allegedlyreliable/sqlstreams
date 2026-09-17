import type { SpawnedLine, SpawnedSquare } from './types';

const idleDelayMs = 5_000;
const growthDurationMs = 37_500;
const spawnAfterMs = 22_500;
const spawnIntervalMs = 3_000;
const spawnLimit = 40;
const spawnFontSizeMin = 18;
const spawnFontSizeMax = 96;
// px per second per px of font size: larger text moves faster
const spawnSpeedPerFontSize = 3;
// a lower bound on the font's average advance, so a copy is never narrower
// than the strip it must span
const spawnMinimumEmPerCharacter = 0.3;
const squareAfterMs = 45_000;
const squareGapMinMs = 500;
const squareGapMaxMs = 2_000;
const squareCountMax = 2;
const squareSizeMin = 20;
const squareSizeMax = 320;
const squareLifeMinMs = 100;
const squareLifeMaxMs = 1_500;
const flickerAfterMs = 45_000;
const flickerRampMs = 67_500;
const flickerGapMaxMs = 4_500;
// three flashes per second is the general flash threshold for
// photosensitive seizures; this floor, even with the gap's jitter, stays
// under two
const flickerGapMinMs = 550;
// long enough for the cover's fade-in to reach its full opacity
const flickerDurationMs = 160;

export class MemberPersonalTextState {
	progress = $state(0);
	scale = $derived(1 + (100 / 13) * this.progress ** 2);
	leftInset = $state(0);
	rightInset = $state(0);
	spawned = $state<SpawnedLine[]>([]);
	squares = $state<SpawnedSquare[]>([]);
	flashing = $state(false);
	held = $state(false);
	private readonly personalText: string;
	private nextSpawnId = 0;
	private nextSquareId = 0;

	constructor(personalText: string) {
		this.personalText = personalText;
	}

	start(text: HTMLElement, textWindow: HTMLElement): () => void {
		const listeners = new AbortController();
		const motion = window.matchMedia('(prefers-reduced-motion: reduce)');
		const bounds = new ResizeObserver(() => {
			const rectangle = textWindow.getBoundingClientRect();
			this.leftInset = rectangle.left;
			this.rightInset = document.documentElement.clientWidth - rectangle.right;
		});
		bounds.observe(textWindow);
		bounds.observe(document.documentElement);
		let timer = 0;
		let spawnTimer = 0;
		let squareTimer = 0;
		let squareRemovals: number[] = [];
		let flickerTimer = 0;
		let flashOffTimer = 0;
		let flickerStarted = 0;
		let frame = 0;
		let started = 0;

		const advance = () => {
			this.setProgress(Math.min(1, (performance.now() - started) / growthDurationMs), text);
			if (this.progress < 1) frame = requestAnimationFrame(advance);
		};

		const spawnNext = () => {
			this.spawn();
			if (this.spawned.length < spawnLimit) {
				spawnTimer = window.setTimeout(spawnNext, spawnIntervalMs);
			}
		};

		// a batch of squares lands together, and each one leaves on its own clock
		const spawnSquares = () => {
			const count = 1 + Math.floor(squareCountMax * Math.random());
			for (let i = 0; i < count; i++) {
				const square = this.spawnSquare();
				const removal = window.setTimeout(
					() => {
						squareRemovals = squareRemovals.filter((other) => other !== removal);
						this.squares = this.squares.filter((existing) => existing.id !== square.id);
					},
					between(squareLifeMinMs, squareLifeMaxMs),
				);
				squareRemovals.push(removal);
			}
			squareTimer = window.setTimeout(spawnSquares, between(squareGapMinMs, squareGapMaxMs));
		};

		const clearSquares = () => {
			window.clearTimeout(squareTimer);
			for (const removal of squareRemovals) window.clearTimeout(removal);
			squareRemovals = [];
			this.squares = [];
		};

		// each flash is one short cover; the gap shrinks along the ramp until
		// the cover stays, then the spawned lines have nothing left to cross
		const flickerNext = () => {
			const ramp = (performance.now() - flickerStarted) / flickerRampMs;
			if (ramp >= 1) {
				window.clearTimeout(spawnTimer);
				this.spawned = [];
				clearSquares();
				this.held = true;
				return;
			}

			this.flashing = true;
			flashOffTimer = window.setTimeout(() => {
				this.flashing = false;
			}, flickerDurationMs);
			const gap = flickerGapMaxMs + (flickerGapMinMs - flickerGapMaxMs) * ramp;
			flickerTimer = window.setTimeout(flickerNext, gap * (0.8 + 0.4 * Math.random()));
		};

		const reset = () => {
			window.clearTimeout(timer);
			window.clearTimeout(spawnTimer);
			window.clearTimeout(flickerTimer);
			window.clearTimeout(flashOffTimer);
			cancelAnimationFrame(frame);
			this.setProgress(0, text);
			this.spawned = [];
			clearSquares();
			this.flashing = false;
			this.held = false;
			if (listeners.signal.aborted || document.hidden || motion.matches) return;

			timer = window.setTimeout(() => {
				started = performance.now();
				frame = requestAnimationFrame(advance);
				spawnTimer = window.setTimeout(spawnNext, spawnAfterMs);
				squareTimer = window.setTimeout(spawnSquares, squareAfterMs);
				flickerTimer = window.setTimeout(() => {
					flickerStarted = performance.now();
					flickerNext();
				}, flickerAfterMs);
			}, idleDelayMs);
		};

		const stop = () => {
			bounds.disconnect();
			listeners.abort();
			reset();
		};

		// Input counts as activity; layout-induced scroll events do not.
		for (const event of [
			'pointermove',
			'pointerdown',
			'keydown',
			'wheel',
			'touchmove',
			'focusin',
		]) {
			window.addEventListener(event, reset, {
				capture: true,
				passive: true,
				signal: listeners.signal,
			});
		}
		document.addEventListener('visibilitychange', reset, { signal: listeners.signal });
		window.addEventListener('pageshow', reset, { signal: listeners.signal });
		motion.addEventListener('change', reset, { signal: listeners.signal });
		document.addEventListener('astro:before-swap', stop, { signal: listeners.signal });
		reset();

		return stop;
	}

	private setProgress(progress: number, text: HTMLElement): void {
		this.progress = progress;
		// Compensate for enlargement while progressively reducing pixel speed.
		const playbackRate = 1 / this.scale ** 1.2;
		for (const animation of text.getAnimations({ subtree: true })) {
			animation.updatePlaybackRate(playbackRate);
		}

		console.log('the end is never the end');
		if (progress >= 1.0) {
			console.log('the end.');
		}
	}

	private spawn(): void {
		const width = document.documentElement.clientWidth;
		const height = document.documentElement.clientHeight;
		const x = width * (0.1 + 0.8 * Math.random());
		const y = height * (0.1 + 0.8 * Math.random());
		const angle = 360 * Math.random();
		const radians = (angle * Math.PI) / 180;
		const fontSize = between(spawnFontSizeMin, spawnFontSizeMax);
		// the strip's ends reach one font size past each edge, so a corner of
		// the strip never shows before its text does
		const entry = reach(x, y, -Math.cos(radians), -Math.sin(radians), width, height) + fontSize;
		const exit = reach(x, y, Math.cos(radians), Math.sin(radians), width, height) + fontSize;
		const distance = entry + exit;
		const phrase = `${this.personalText} `;
		// two copies loop by one copy's width, so one copy must span the strip
		const repeats =
			Math.ceil(distance / (phrase.length * fontSize * spawnMinimumEmPerCharacter)) + 1;
		this.spawned.push({
			id: this.nextSpawnId++,
			x,
			y,
			angle,
			fontSize,
			speed: fontSize * spawnSpeedPerFontSize,
			entry,
			distance,
			copy: phrase.repeat(repeats),
			copyWidth: 0,
		});
	}

	private spawnSquare(): SpawnedSquare {
		const size = between(squareSizeMin, squareSizeMax);
		const square = {
			id: this.nextSquareId++,
			x: between(0, document.documentElement.clientWidth - size),
			y: between(0, document.documentElement.clientHeight - size),
			size,
		};
		this.squares.push(square);
		return square;
	}
}

// ***************
// *** HELPERS ***
// ***************

function between(min: number, max: number): number {
	return min + (max - min) * Math.random();
}

// the distance from a viewport point to the viewport edge along a direction
function reach(
	x: number,
	y: number,
	directionX: number,
	directionY: number,
	width: number,
	height: number,
): number {
	const alongX =
		directionX > 0 ? (width - x) / directionX : directionX < 0 ? -x / directionX : Infinity;
	const alongY =
		directionY > 0 ? (height - y) / directionY : directionY < 0 ? -y / directionY : Infinity;
	return Math.min(alongX, alongY);
}
