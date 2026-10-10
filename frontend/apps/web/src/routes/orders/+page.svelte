<script lang="ts">
	import { Package } from '@lucide/svelte';
	import { page } from '$app/state';
	import { formatDate, formatPrice } from '@ecommerce/lib/format';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import { EmptyState } from '@ecommerce/ui/empty-state';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import { getOrders } from './orders.remote';

	const currentPage = $derived(Math.max(0, Number(page.url.searchParams.get('page') ?? 0) || 0));
	const args = $derived({ page: currentPage, size: 10 });

	const href = (target: number) =>
		localizeHref(target === 0 ? '/orders' : `/orders?page=${target}`);
</script>

<svelte:head>
	<title>{m.orders_metaTitle()}</title>
</svelte:head>

<h1 class="mb-4 text-2xl font-bold">{m.orders_title()}</h1>

{#await getOrders(args)}
	<div class="grid gap-3">
		{#each Array(3).keys() as i (i)}
			<Skeleton class="h-28 w-full rounded-xl" />
		{/each}
	</div>
{:then result}
	{#if result.content.length === 0}
		<EmptyState icon={Package} title={m.orders_emptyTitle()} description={m.orders_emptyHint()}>
			{#snippet action()}
				<Button href={localizeHref('/products')}>{m.orders_shopNow()}</Button>
			{/snippet}
		</EmptyState>
	{:else}
		<div class="grid gap-3">
			{#each result.content as order (order.orderId)}
				<div class="rounded-xl border border-border bg-card p-4">
					<div class="flex flex-wrap items-center gap-3">
						<Badge variant="secondary">{m.orders_processing()}</Badge>
						<span class="font-semibold">{m.orders_orderLabel({ id: String(order.orderId) })}</span>
						<span class="text-sm text-muted-foreground">{formatDate(order.orderDate)}</span>
						<span class="ml-auto font-bold text-primary">{formatPrice(order.orderFee)}</span>
					</div>
					{#if order.orderDesc}
						<p class="mt-2 text-sm text-muted-foreground">{order.orderDesc}</p>
					{/if}
				</div>
			{/each}
		</div>

		{#if result.totalPages > 1}
			<div class="mt-6 flex items-center justify-center gap-2">
				<Button
					variant="outline"
					size="sm"
					href={href(currentPage - 1)}
					disabled={currentPage <= 0}
				>
					{m.common_prev()}
				</Button>
				<span class="text-sm text-muted-foreground">
					{m.common_page({ current: String(currentPage + 1), total: String(result.totalPages) })}
				</span>
				<Button
					variant="outline"
					size="sm"
					href={href(currentPage + 1)}
					disabled={currentPage >= result.totalPages - 1}
				>
					{m.common_next()}
				</Button>
			</div>
		{/if}
	{/if}
{:catch}
	<div
		class="grid justify-items-center gap-3 rounded-xl border border-border bg-card p-12 text-center"
	>
		<p class="font-semibold">{m.common_error()}</p>
		<Button href={localizeHref('/orders')} variant="outline" size="sm">{m.common_retry()}</Button>
	</div>
{/await}
