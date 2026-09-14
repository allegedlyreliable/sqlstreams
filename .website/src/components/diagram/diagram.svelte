<script lang="ts">
	import { MarkerType, SvelteFlow } from '@xyflow/svelte';
	import ResourceCard from '../resource-card/resource-card.svelte';
	import type { DiagramLayout } from './types';

	type Props = { graph: DiagramLayout; description: string; caption: string };
	let { graph, description, caption }: Props = $props();
	const nodeTypes = { resource: ResourceCard };
	const zoom = 1.15;
	const defaultEdgeOptions = {
		type: 'smoothstep',
		animated: false,
		markerEnd: { type: MarkerType.ArrowClosed, color: 'var(--color-diagram-arrow)' },
		style: 'stroke: var(--color-diagram-arrow); stroke-dasharray: none; animation: none',
	};
</script>

<figure class="diagram">
	<div class="scroll">
		<div
			class="canvas"
			style:width={`${graph.width * zoom}px`}
			style:height={`${graph.height * zoom}px`}
		>
			<div class="drawing" role="img" aria-label={description}>
				<SvelteFlow
					nodes={graph.nodes}
					edges={graph.edges}
					{nodeTypes}
					{defaultEdgeOptions}
					width={graph.width * zoom}
					height={graph.height * zoom}
					viewport={{ x: 0, y: 0, zoom }}
					nodesDraggable={false}
					nodesConnectable={false}
					nodesFocusable={false}
					edgesFocusable={false}
					elementsSelectable={false}
					disableKeyboardA11y
					panOnDrag={false}
					zoomOnScroll={false}
					zoomOnPinch={false}
					zoomOnDoubleClick={false}
					preventScrolling={false}
					proOptions={{ hideAttribution: true }}
				/>
			</div>
			<!-- praise the open sourcers -->
			<a class="attribution" href="https://svelteflow.dev" target="_blank" rel="noopener noreferrer"
				>Svelte Flow</a
			>
		</div>
	</div>
	<figcaption>{caption}</figcaption>
</figure>

<style src="./diagram.css"></style>
