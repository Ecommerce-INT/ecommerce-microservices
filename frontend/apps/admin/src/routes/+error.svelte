<script lang="ts">
	import { AlertTriangle, SearchX } from '@lucide/svelte';
	import { page } from '$app/state';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import { EmptyState } from '@ecommerce/ui/empty-state';

	const notFound = $derived(page.status === 404);
</script>

<svelte:head>
	<title>{m.error_metaTitle({ status: String(page.status) })}</title>
</svelte:head>

<div class="mx-auto mt-16 max-w-lg">
	<EmptyState
		icon={notFound ? SearchX : AlertTriangle}
		title={notFound ? m.error_notFoundTitle() : m.error_title()}
		description={notFound ? m.error_notFoundHint() : m.error_hint()}
	>
		{#snippet action()}
			<div class="flex flex-wrap items-center justify-center gap-2">
				<Button href={localizeHref('/')}>{m.error_backHome()}</Button>
				<Button variant="outline" onclick={() => location.reload()}>{m.error_retry()}</Button>
			</div>
		{/snippet}
	</EmptyState>
</div>
