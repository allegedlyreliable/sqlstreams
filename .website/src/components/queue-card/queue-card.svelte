<script lang="ts">
	import type { QueueSnapshot } from './types';
	type Props = { data: QueueSnapshot };
	let { data }: Props = $props();
	const waiting = $derived([...data.waiting].reverse());
</script>

<div class="snapshot">
	<div class="stream">
		<header><strong>Stream: {data.stream}</strong></header>
		<div class="stored">Stored messages: {[data.processing, ...data.waiting].join(' · ')}</div>
	</div>
	<div class="claim">
		<span class="direction" aria-hidden="true">↓</span>
		<div>
			<strong
				>Consumer instance {data.instance} claims these {data.waiting.length + 1} messages</strong
			>
			<div>Lease clock starts at this claim</div>
		</div>
	</div>
	<div class="group">
		<header><strong>Consumer group: {data.group}</strong></header>
		<div class="group-body">
			<div class="consumer">
				<header><strong>Consumer instance {data.instance}</strong></header>
				<div class="work">
					<div class="label">Waiting in {data.instance}'s queue</div>
					<div class="label processing-label">Processing</div>
					<div class="waiting">
						{#each waiting as id (id)}<div class="tile">{id}</div>{/each}
					</div>
					<div class="direction" aria-hidden="true">→</div>
					<div class="tile" data-active="true">{data.processing}</div>
					<div class="lease">
						<strong>One lease covers all {data.waiting.length + 1} messages</strong>
						<span>Counts down while messages wait</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style src="./queue-card.css"></style>
