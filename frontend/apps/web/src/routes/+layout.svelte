<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Toaster } from '@ecommerce/ui/sonner';
	import Footer from '#lib/components/footer.svelte';
	import Header from '#lib/components/header.svelte';
	import { cart } from '#lib/stores/cart.svelte';
	import { baseLocale, locales, localizeHref } from '#lib/paraglide/runtime.js';
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	onMount(() => cart.hydrate());

	const languages = $derived(
		locales.map((locale) => ({
			locale,
			href: new URL(localizeHref(page.url.pathname, { locale }), page.url.origin).href
		}))
	);
	const defaultHref = $derived(languages.find((language) => language.locale === baseLocale)?.href);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="canonical" href={page.url.href} />
	{#each languages as language (language.locale)}
		<link rel="alternate" hreflang={language.locale} href={language.href} />
	{/each}
	{#if defaultHref}
		<link rel="alternate" hreflang="x-default" href={defaultHref} />
	{/if}
</svelte:head>

<div class="flex min-h-screen flex-col">
	<Header user={data.user} />
	<main class="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 lg:px-8">
		{@render children()}
	</main>
	<Footer />
</div>

<Toaster richColors position="top-center" />
