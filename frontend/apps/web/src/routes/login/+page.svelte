<script lang="ts">
	import { ShieldCheck } from '@lucide/svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { Button } from '@ecommerce/ui/button';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{m.auth_loginPageTitle()}</title>
</svelte:head>

<div class="mx-auto grid max-w-md gap-4 py-10">
	<div class="rounded-xl border border-border bg-card p-6">
		<h1 class="text-xl font-bold">{m.auth_loginTitle()}</h1>
		<p class="mt-1 text-sm text-muted-foreground">{m.auth_loginSubtitle()}</p>

		<Button href={data.ssoUrl} size="lg" class="mt-5 w-full">
			{m.auth_ssoLoginBtn()}
		</Button>
		<p class="mt-2 text-center text-xs text-muted-foreground">{m.auth_ssoRedirectNote()}</p>

		<div class="mt-6 rounded-lg bg-muted p-4">
			<p class="flex items-center gap-2 text-sm font-semibold">
				<ShieldCheck class="size-4 text-primary" />
				{m.auth_ssoWhy()}
			</p>
			<ul class="mt-2 grid gap-1 text-xs text-muted-foreground">
				<li>• {m.auth_ssoFeature1()}</li>
				<li>• {m.auth_ssoFeature2()}</li>
				<li>• {m.auth_ssoFeature3()}</li>
			</ul>
		</div>

		<div class="mt-6 border-t border-border pt-4">
			<p class="text-xs text-muted-foreground">{m.auth_registerSsoNote()}</p>
			<Button href={data.ssoUrl} variant="outline" class="mt-2 w-full">
				{m.auth_registerSsoBtn()}
			</Button>
		</div>
	</div>

	<p class="text-center text-xs text-muted-foreground">{m.auth_copyright()}</p>
</div>
