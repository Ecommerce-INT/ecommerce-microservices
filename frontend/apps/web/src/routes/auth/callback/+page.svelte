<script lang="ts">
	import { CircleAlert, Loader2 } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{m.auth_loginPageTitle()}</title>
</svelte:head>

<div class="mx-auto grid max-w-md gap-3 py-16 text-center">
	{#if data.failed}
		<CircleAlert class="mx-auto size-10 text-destructive" />
		<h1 class="text-lg font-semibold">{m.auth_callbackErrorTitle()}</h1>
		<p class="text-sm text-muted-foreground">{m.auth_callbackErrorText()}</p>
		<Button href={localizeHref('/login')} class="mx-auto mt-2">
			{m.auth_callbackRetry()}
		</Button>
	{:else}
		<Loader2 class="mx-auto size-10 animate-spin text-primary" />
		<p class="text-sm text-muted-foreground">{m.auth_callbackLoading()}</p>
	{/if}
</div>
