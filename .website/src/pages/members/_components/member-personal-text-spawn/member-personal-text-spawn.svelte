<script lang="ts">
	import type { SpawnedLine } from '../member-personal-text/types';

	type Props = {
		line: SpawnedLine;
		// receives one copy's rendered width whenever it changes
		onCopyWidth: (width: number) => void;
	};

	let { line, onCopyWidth }: Props = $props();

	// the loop shifts by one copy's width, so its duration needs the real width
	function measureCopy(copy: HTMLElement): () => void {
		const bounds = new ResizeObserver(() => onCopyWidth(copy.offsetWidth));
		bounds.observe(copy);
		return () => bounds.disconnect();
	}
</script>

<div
	class="personal-text-spawn"
	style:transform={`translate(${line.x}px, ${line.y}px) rotate(${line.angle}deg)`}
	style:font-size={`${line.fontSize}px`}
	style:--personal-text-spawn-entry={`${line.entry}px`}
	style:--personal-text-spawn-distance={`${line.distance}px`}
	style:--personal-text-spawn-roll-duration={`${line.distance / line.speed}s`}
	style:--personal-text-spawn-loop-duration={`${line.copyWidth / line.speed}s`}
>
	<div class="personal-text-spawn-track">
		<span class="personal-text-spawn-copy" {@attach measureCopy}>{line.copy}</span>
		<span class="personal-text-spawn-copy">{line.copy}</span>
	</div>
</div>

<style src="./member-personal-text-spawn.css"></style>
