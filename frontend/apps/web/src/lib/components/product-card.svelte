<script lang="ts">
	import { ShoppingCart, Star } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { formatPrice } from '@ecommerce/lib/format';
	import type { Product } from '@ecommerce/lib/types';
	import { m } from '#lib/paraglide/messages.js';
	import { cart } from '#lib/stores/cart.svelte';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import type { ShowcaseProduct } from '#lib/data/home-mock';

	type CardProduct = Product | ShowcaseProduct;

	let { product, class: className = '' }: { product: CardProduct; class?: string } = $props();

	const stock = $derived('quantity' in product ? product.quantity : undefined);
	const outOfStock = $derived(stock !== undefined && stock <= 0);

	function addToCart(event: MouseEvent) {
		event.preventDefault();
		cart.add(
			{
				productId: product.productId,
				productTitle: product.productTitle,
				imageUrl: product.imageUrl,
				priceUnit: product.priceUnit,
				quantity: stock ?? 1,
				category: product.category
			},
			1
		);
		toast.success(m.productcard_added());
	}
</script>

<a
	href="/products/{product.productId}"
	class="group flex flex-col overflow-hidden rounded-xl border border-border bg-card transition-shadow hover:shadow-md {className}"
>
	<div class="relative aspect-square overflow-hidden bg-muted">
		{#if product.imageUrl}
			<img
				src={product.imageUrl}
				alt={product.productTitle}
				class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
				loading="lazy"
			/>
		{:else}
			<div class="flex h-full w-full items-center justify-center text-5xl text-muted-foreground">
				📦
			</div>
		{/if}

		{#if 'discountPercent' in product && product.discountPercent}
			<span
				class="absolute top-2 left-2 rounded-md bg-primary px-1.5 py-0.5 text-xs font-semibold text-primary-foreground"
			>
				-{product.discountPercent}%
			</span>
		{:else if outOfStock}
			<span class="absolute top-2 left-2">
				<Badge variant="secondary">{m.productcard_outOfStock()}</Badge>
			</span>
		{/if}

		{#if 'badge' in product && product.badge}
			<span class="absolute bottom-2 left-2">
				<Badge>{product.badge}</Badge>
			</span>
		{/if}
	</div>

	<div class="flex flex-1 flex-col gap-1.5 p-3">
		<h3 class="line-clamp-2 min-h-10 text-sm leading-5 font-medium">{product.productTitle}</h3>

		<div class="flex items-baseline gap-2 text-primary">
			<span class="text-base font-bold">{formatPrice(product.priceUnit)}</span>
			{#if 'oldPrice' in product && product.oldPrice}
				<span class="text-xs text-muted-foreground line-through"
					>{formatPrice(product.oldPrice)}</span
				>
			{/if}
		</div>

		<div class="flex items-center gap-2 text-xs text-muted-foreground">
			{#if 'rating' in product && product.rating}
				<span class="flex items-center gap-0.5 text-amber-500">
					{#each Array(product.rating).keys() as i (i)}
						<Star class="size-3 fill-current" />
					{/each}
				</span>
			{/if}
			{#if 'sold' in product && product.sold}
				<span>{m.home_sold({ count: String(product.sold) })}</span>
			{/if}
			{#if stock !== undefined && stock > 0 && stock <= 10}
				<span>{m.productcard_remaining({ count: String(stock) })}</span>
			{/if}
		</div>

		<div class="mt-auto pt-2">
			<Button size="sm" class="w-full" disabled={outOfStock} onclick={addToCart}>
				<ShoppingCart class="size-4" />
				{m.productcard_addToCart()}
			</Button>
		</div>
	</div>
</a>
