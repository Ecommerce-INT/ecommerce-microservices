<script lang="ts">
	import {
		ChevronRight,
		Minus,
		Plus,
		RefreshCw,
		ShieldCheck,
		ShoppingCart,
		Truck
	} from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import { formatPrice } from '@ecommerce/lib/format';
	import { m } from '$lib/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { cart } from '$lib/stores/cart.svelte';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import { getProduct } from './product.remote';

	const productId = $derived(Number(page.params.id));
	let quantity = $state(1);
</script>

<div class="mb-4 flex items-center gap-1 text-sm">
	<a href={localizeHref('/')} class="text-muted-foreground hover:text-primary">
		{m.productdetail_breadcrumbHome()}
	</a>
	<ChevronRight class="size-3.5 text-muted-foreground" />
	<a href={localizeHref('/products')} class="text-muted-foreground hover:text-primary">
		{m.productdetail_breadcrumbProducts()}
	</a>
</div>

{#await getProduct(productId)}
	<div class="grid gap-6 lg:grid-cols-2">
		<Skeleton class="aspect-square w-full rounded-xl" />
		<div class="grid content-start gap-3">
			<Skeleton class="h-8 w-3/4" />
			<Skeleton class="h-5 w-1/3" />
			<Skeleton class="h-10 w-1/2" />
			<Skeleton class="h-9 w-40" />
		</div>
	</div>
{:then product}
	<div class="grid gap-8 lg:grid-cols-2">
		<div class="overflow-hidden rounded-xl border border-border bg-card">
			{#if product.imageUrl}
				<img
					src={product.imageUrl}
					alt={product.productTitle}
					class="aspect-square w-full object-cover"
				/>
			{:else}
				<div
					class="flex aspect-square w-full items-center justify-center text-7xl text-muted-foreground"
				>
					📦
				</div>
			{/if}
		</div>

		<div class="grid content-start gap-4">
			<h1 class="text-2xl font-bold">{product.productTitle}</h1>

			{#if product.sku}
				<p class="text-sm text-muted-foreground">SKU: {product.sku}</p>
			{/if}

			<div class="rounded-xl bg-muted/50 p-4">
				<div class="flex items-baseline gap-3">
					<span class="text-3xl font-bold text-primary">{formatPrice(product.priceUnit)}</span>
					<span class="text-sm text-muted-foreground">{m.productdetail_vatIncluded()}</span>
				</div>
			</div>

			<div class="flex items-center gap-3">
				{#if product.quantity > 0}
					<Badge variant="secondary"
						>{m.productdetail_inStock({ count: String(product.quantity) })}</Badge
					>
				{:else}
					<Badge variant="destructive">{m.productdetail_outOfStock()}</Badge>
				{/if}
			</div>

			<div class="flex items-center gap-4">
				<span class="text-sm font-medium">{m.productdetail_quantityLabel()}</span>
				<div class="flex items-center rounded-lg border border-border">
					<Button
						variant="ghost"
						size="icon-sm"
						aria-label="decrease"
						onclick={() => (quantity = Math.max(1, quantity - 1))}
					>
						<Minus class="size-3.5" />
					</Button>
					<span class="w-10 text-center text-sm font-semibold">{quantity}</span>
					<Button
						variant="ghost"
						size="icon-sm"
						aria-label="increase"
						onclick={() => (quantity = Math.min(product.quantity || 99, quantity + 1))}
					>
						<Plus class="size-3.5" />
					</Button>
				</div>
			</div>

			<div class="grid gap-3 sm:grid-cols-2">
				<Button
					size="lg"
					disabled={product.quantity <= 0}
					onclick={() => {
						cart.add(product, quantity);
						toast.success(m.productdetail_added());
					}}
				>
					<ShoppingCart class="size-4" />
					{m.productdetail_addToCart()}
				</Button>
				<Button
					size="lg"
					variant="outline"
					disabled={product.quantity <= 0}
					href={localizeHref('/checkout')}
					onclick={() => {
						cart.add(product, quantity);
					}}
				>
					{m.productdetail_buyNow()}
				</Button>
			</div>

			<ul class="grid gap-2 text-sm text-muted-foreground">
				<li class="flex items-center gap-2">
					<Truck class="size-4 text-primary" />
					{m.productdetail_freeShipping()}
				</li>
				<li class="flex items-center gap-2">
					<ShieldCheck class="size-4 text-primary" />
					{m.productdetail_warranty()}
				</li>
				<li class="flex items-center gap-2">
					<RefreshCw class="size-4 text-primary" />
					{m.productdetail_freeReturn()}
				</li>
			</ul>

			{#if product.description}
				<div class="border-t border-border pt-4">
					<h2 class="mb-2 font-semibold">{m.productdetail_description()}</h2>
					<p class="text-sm whitespace-pre-line text-muted-foreground">{product.description}</p>
				</div>
			{/if}
		</div>
	</div>
{:catch}
	<div
		class="grid justify-items-center gap-3 rounded-xl border border-border bg-card p-12 text-center"
	>
		<p class="text-lg font-semibold">{m.productdetail_notFound()}</p>
		<Button href={localizeHref('/products')} variant="outline">
			{m.productdetail_backToList()}
		</Button>
	</div>
{/await}
