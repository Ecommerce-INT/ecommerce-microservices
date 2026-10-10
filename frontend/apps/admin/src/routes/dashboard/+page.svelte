<script lang="ts">
	import { Package, ShoppingBag, Tags } from '@lucide/svelte';
	import { formatDate, formatPrice } from '@ecommerce/lib/format';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Badge } from '@ecommerce/ui/badge';
	import * as Table from '@ecommerce/ui/table';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import { getDashboard } from './dashboard.remote';

	const stats = [
		{ key: 'products', icon: Package, label: m.dashboard_products, href: '/products' },
		{ key: 'orders', icon: ShoppingBag, label: m.dashboard_orders, href: '/orders' },
		{ key: 'categories', icon: Tags, label: m.sidebar_categories, href: '/categories' }
	] as const;
</script>

<svelte:head>
	<title>{m.meta_template({ page: m.dashboard_metaTitle() })}</title>
</svelte:head>

<div class="mb-5">
	<h1 class="text-xl font-bold">{m.dashboard_title()}</h1>
	<p class="text-sm text-muted-foreground">{m.dashboard_subtitle()}</p>
</div>

{#await getDashboard()}
	<div class="grid gap-4 sm:grid-cols-3">
		{#each stats as stat (stat.key)}
			<Skeleton class="h-24 rounded-xl" />
		{/each}
	</div>
	<Skeleton class="mt-6 h-64 rounded-xl" />
{:then data}
	<div class="grid gap-4 sm:grid-cols-3">
		{#each stats as stat (stat.key)}
			{@const value =
				stat.key === 'products'
					? data.productCount
					: stat.key === 'orders'
						? data.orderCount
						: data.categoryCount}
			<a
				href={localizeHref(stat.href)}
				class="rounded-xl border border-border bg-card p-5 transition-colors hover:border-primary"
			>
				<div class="flex items-center justify-between">
					<span class="text-sm text-muted-foreground">{stat.label()}</span>
					<stat.icon class="size-5 text-primary" />
				</div>
				<p class="mt-2 text-3xl font-bold">{value}</p>
			</a>
		{/each}
	</div>

	<div class="mt-6 grid gap-4 xl:grid-cols-2">
		<section class="rounded-xl border border-border bg-card">
			<header class="flex items-center justify-between border-b border-border px-5 py-3">
				<h2 class="font-semibold">{m.dashboard_recentOrders()}</h2>
				<a href={localizeHref('/orders')} class="text-sm text-primary">
					{m.common_viewAll()}
				</a>
			</header>
			{#if data.recentOrders.length === 0}
				<p class="p-5 text-sm text-muted-foreground">{m.dashboard_noOrders()}</p>
			{:else}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.orders_colId()}</Table.Head>
							<Table.Head>{m.orders_colDate()}</Table.Head>
							<Table.Head class="text-right">{m.orders_colTotal()}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each data.recentOrders as order (order.orderId)}
							<Table.Row>
								<Table.Cell class="font-medium">
									{m.dashboard_order({ id: String(order.orderId) })}
								</Table.Cell>
								<Table.Cell class="text-muted-foreground">{formatDate(order.orderDate)}</Table.Cell>
								<Table.Cell class="text-right font-medium">{formatPrice(order.orderFee)}</Table.Cell
								>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			{/if}
		</section>

		<section class="rounded-xl border border-border bg-card">
			<header class="flex items-center justify-between border-b border-border px-5 py-3">
				<h2 class="font-semibold">{m.dashboard_recentProducts()}</h2>
				<a href={localizeHref('/products')} class="text-sm text-primary">
					{m.common_viewAll()}
				</a>
			</header>
			{#if data.recentProducts.length === 0}
				<p class="p-5 text-sm text-muted-foreground">{m.dashboard_noProducts()}</p>
			{:else}
				<ul class="divide-y divide-border">
					{#each data.recentProducts as product (product.productId)}
						<li class="flex items-center gap-3 px-5 py-3">
							<div class="size-10 shrink-0 overflow-hidden rounded-lg bg-muted">
								{#if product.imageUrl}
									<img
										src={product.imageUrl}
										alt={product.productTitle}
										class="h-full w-full object-cover"
									/>
								{/if}
							</div>
							<div class="min-w-0 flex-1">
								<p class="truncate text-sm font-medium">{product.productTitle}</p>
								<p class="text-xs text-muted-foreground">
									{product.category?.categoryTitle ?? m.dashboard_uncategorized()}
								</p>
							</div>
							<div class="text-right">
								<p class="text-sm font-medium">{formatPrice(product.priceUnit)}</p>
								{#if product.quantity > 0}
									<Badge variant="secondary">
										{m.dashboard_inStock({ count: String(product.quantity) })}
									</Badge>
								{:else}
									<Badge variant="destructive">{m.dashboard_outOfStock()}</Badge>
								{/if}
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</section>
	</div>
{:catch}
	<p class="text-sm text-destructive">{m.common_noData()}</p>
{/await}
