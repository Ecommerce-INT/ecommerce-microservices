<script lang="ts">
	import { onMount } from 'svelte';
	import { Toaster } from '@ecommerce/ui/sonner';
	import Footer from '$lib/components/footer.svelte';
	import Header from '$lib/components/header.svelte';
	import { cart } from '$lib/stores/cart.svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	onMount(() => cart.hydrate());
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex min-h-screen flex-col">
	<Header user={data.user} />
	<main class="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 lg:px-8">
		{@render children()}
	</main>
	<Footer />
</div>

<Toaster richColors position="top-center" />
