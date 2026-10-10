<script lang="ts">
	import { CircleAlert, Loader2 } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{m.meta_title()}</title>
</svelte:head>

<div class="grid min-h-screen place-items-center p-4">
	<div class="grid justify-items-center gap-3 text-center">
		{#if data.failed}
			<CircleAlert class="size-10 text-destructive" />
			<h1 class="text-lg font-semibold">{m.auth_callbackErrorTitle()}</h1>
			<p class="text-sm text-muted-foreground">{m.auth_callbackErrorText()}</p>
			<Button href={localizeHref('/login')} class="mt-2">{m.auth_callbackRetry()}</Button>
		{:else}
			<Loader2 class="size-10 animate-spin text-primary" />
			<p class="text-sm text-muted-foreground">{m.auth_callbackLoading()}</p>
		{/if}
	</div>
</div>
