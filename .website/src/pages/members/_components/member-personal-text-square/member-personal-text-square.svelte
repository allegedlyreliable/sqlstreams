<script lang="ts">
	import type { SpawnedSquare } from '../member-personal-text/types';

	type Props = {
		square: SpawnedSquare;
	};

	let { square }: Props = $props();
</script>

<div
	class="personal-text-square"
	style:transform={`translate(${square.x}px, ${square.y}px)`}
	style:width={`${square.size}px`}
	style:height={`${square.size}px`}
></div>

<style src="./member-personal-text-square.css"></style>
