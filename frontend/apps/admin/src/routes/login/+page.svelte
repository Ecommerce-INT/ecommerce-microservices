<script lang="ts">
	import { ShieldCheck } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { Button } from '@ecommerce/ui/button';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{m.meta_title()}</title>
</svelte:head>

<div class="grid min-h-screen place-items-center p-4">
	<div class="w-full max-w-md rounded-xl border border-border bg-card p-6">
		<p class="text-2xl font-black text-primary">{m.sidebar_brand()}</p>
		<h1 class="mt-3 text-lg font-semibold">{m.login_title()}</h1>
		<p class="mt-1 text-sm text-muted-foreground">{m.login_note()}</p>

		<Button href={data.ssoUrl} size="lg" class="mt-5 w-full">
			<ShieldCheck class="size-4" />
			{m.login_btn()}
		</Button>
	</div>
</div>
