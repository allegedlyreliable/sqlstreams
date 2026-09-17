<script lang="ts">
	import { boards, boardHref } from '../../pages/_board/navigation';
	import { repositoryUrl } from '../../site';

	type NavLink = {
		label: string;
		href: string;
		inPhoneNav: boolean;
	};

	type Props = {
		// null on a page no nav link names (a thread, the 404 page)
		current: string | null;
	};

	let { current }: Props = $props();

	const links: NavLink[] = [
		{ label: 'Board Index', href: '/', inPhoneNav: false },
		...boards.map((board) => ({
			label: board.title,
			href: boardHref(board),
			inPhoneNav: board.inPhoneNav,
		})),
		{ label: 'Search', href: '/search/', inPhoneNav: true },
		{ label: 'GitHub', href: repositoryUrl, inPhoneNav: true },
	];
</script>

<nav class="board-nav" aria-label="Board">
	{#each links as link (link.label)}
		<a
			href={link.href}
			aria-current={link.label === current ? 'page' : undefined}
			data-phone-hidden={link.inPhoneNav ? undefined : ''}
		>
			{link.label}
		</a>
	{/each}
</nav>

<style src="./board-nav.css"></style>
