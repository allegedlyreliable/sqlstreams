<script lang="ts">
	import { onMount } from 'svelte';
	import { MemberPersonalTextState } from './member-personal-text-state.svelte';
	import MemberPersonalTextSpawn from '../member-personal-text-spawn/member-personal-text-spawn.svelte';
	import MemberPersonalTextSquare from '../member-personal-text-square/member-personal-text-square.svelte';
	import MemberPersonalTextFlash from '../member-personal-text-flash/member-personal-text-flash.svelte';

	type Props = {
		// the phrase the personal text repeats without ever completing
		personalText: string;
		// the line the flash cover finally holds still
		finalText: string;
	};

	let { personalText, finalText }: Props = $props();

	// each copy must be wider than any viewport, so the loop never shows a gap
	const repeated = `${personalText} `.repeat(20);
	const state = new MemberPersonalTextState(personalText);
	let text: HTMLDivElement;
	let textWindow: HTMLDivElement;

	onMount(() => state.start(text, textWindow));
</script>

<div
	class="personal-text"
	data-idle={state.progress > 0}
	style:--personal-text-progress={state.progress}
	style:--personal-text-scale={state.scale}
	style:--personal-text-left-inset={`${state.leftInset}px`}
	style:--personal-text-right-inset={`${state.rightInset}px`}
>
	<div class="personal-text-veil" aria-hidden="true"></div>
	<div class="personal-text-spawn-layer" aria-hidden="true">
		{#each state.spawned as line (line.id)}
			<MemberPersonalTextSpawn
				{line}
				onCopyWidth={(width) => {
					line.copyWidth = width;
				}}
			/>
		{/each}
	</div>
	<div class="personal-text-square-layer" aria-hidden="true">
		{#each state.squares as square (square.id)}
			<MemberPersonalTextSquare {square} />
		{/each}
	</div>
	<MemberPersonalTextFlash flashing={state.flashing} held={state.held} text={finalText} />
	<span class="personal-text-label">Personal text</span>
	<div class="personal-text-window" bind:this={textWindow}>
		<div class="personal-text-growth">
			<div class="personal-text-track" bind:this={text}>
				<span class="personal-text-copy">{repeated}</span>
				<span class="personal-text-copy" aria-hidden="true">{repeated}</span>
			</div>
		</div>
	</div>
</div>

<style src="./member-personal-text.css"></style>
