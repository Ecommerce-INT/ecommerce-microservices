<script lang="ts">
	import { CircleSlash } from '@lucide/svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import { logout } from '$lib/remote/auth.remote';

	let loggingOut = $state(false);

	async function handleLogout() {
		if (loggingOut) return;
		loggingOut = true;
		try {
			await logout();
			await invalidateAll();
			await goto(localizeHref('/login'));
		} finally {
			loggingOut = false;
		}
	}
</script>

<svelte:head>
	<title>{m.access_denied_title()} | {m.meta_title()}</title>
</svelte:head>

<div class="grid min-h-screen place-items-center p-4">
	<div
		class="grid w-full max-w-md justify-items-center gap-3 rounded-xl border border-border bg-card p-8 text-center"
	>
		<CircleSlash class="size-12 text-destructive" />
		<h1 class="text-lg font-semibold">{m.access_denied_title()}</h1>
		<p class="text-sm text-muted-foreground">{m.access_denied_text()}</p>
		<Button variant="outline" class="mt-2" disabled={loggingOut} onclick={handleLogout}>
			{m.sidebar_logout()}
		</Button>
	</div>
</div>
