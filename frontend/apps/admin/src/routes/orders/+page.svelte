<script lang="ts">
	import { Search } from '@lucide/svelte';
	import { page } from '$app/state';
	import { formatDate, formatPrice } from '@ecommerce/lib/format';
	import type { Order } from '@ecommerce/lib/types';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import { Input } from '@ecommerce/ui/input';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import * as Table from '@ecommerce/ui/table';
	import { getOrders } from './orders.remote';

	const currentPage = $derived(Math.max(0, Number(page.url.searchParams.get('page') ?? 0) || 0));
	const args = $derived({ page: currentPage, size: 10 });

	let filter = $state('');

	const href = (target: number) =>
		localizeHref(target === 0 ? '/orders' : `/orders?page=${target}`);

	const visible = (orders: Order[]) =>
		filter.trim()
			? orders.filter((order) => String(order.orderId).includes(filter.trim()))
			: orders;
</script>

<svelte:head>
	<title>{m.meta_template({ page: m.orders_metaTitle() })}</title>
</svelte:head>

<div class="mb-4 flex flex-wrap items-center gap-3">
	<h1 class="text-xl font-bold">{m.orders_title()}</h1>
	<div class="relative ml-auto">
		<Search
			class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground"
		/>
		<Input
			bind:value={filter}
			placeholder={m.orders_searchPlaceholder()}
			aria-label={m.common_search()}
			class="h-9 w-60 pl-8"
		/>
	</div>
</div>

{#await getOrders(args)}
	<Skeleton class="h-72 rounded-xl" />
{:then result}
	{@const rows = visible(result.content)}
	<p class="mb-2 text-sm text-muted-foreground">
		{m.orders_count({ count: String(result.totalElements) })}
	</p>

	<div class="overflow-hidden rounded-xl border border-border bg-card">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>{m.orders_colId()}</Table.Head>
					<Table.Head>{m.orders_colDate()}</Table.Head>
					<Table.Head>{m.orders_colDesc()}</Table.Head>
					<Table.Head class="text-right">{m.orders_colTotal()}</Table.Head>
					<Table.Head class="text-center">{m.orders_colStatus()}</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each rows as order (order.orderId)}
					<Table.Row>
						<Table.Cell class="font-medium"
							>{m.orders_order({ id: String(order.orderId) })}</Table.Cell
						>
						<Table.Cell class="text-muted-foreground">{formatDate(order.orderDate)}</Table.Cell>
						<Table.Cell class="max-w-64 truncate text-muted-foreground">
							{order.orderDesc ?? '—'}
						</Table.Cell>
						<Table.Cell class="text-right font-medium">{formatPrice(order.orderFee)}</Table.Cell>
						<Table.Cell class="text-center">
							<Badge variant="secondary">{m.orders_processing()}</Badge>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>

		{#if rows.length === 0}
			<p class="p-10 text-center text-sm text-muted-foreground">{m.orders_noOrders()}</p>
		{/if}
	</div>

	{#if result.totalPages > 1}
		<div class="mt-4 flex items-center justify-center gap-2">
			<Button variant="outline" size="sm" href={href(currentPage - 1)} disabled={currentPage <= 0}>
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
{:catch}
	<p class="text-sm text-destructive">{m.common_noData()}</p>
{/await}
