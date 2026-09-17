<script lang="ts">
	type Props = {
		// the cover is showing for one flash
		flashing: boolean;
		// the cover stays and shows the text
		held: boolean;
		text: string;
	};

	let { flashing, held, text }: Props = $props();
</script>

<div class="personal-text-flash" data-flashing={flashing} data-held={held} aria-hidden="true">
	{#if held}
		<span class="personal-text-flash-text">{text}</span>
	{/if}
</div>

<style src="./member-personal-text-flash.css"></style>
