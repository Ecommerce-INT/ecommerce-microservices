<script lang="ts">
	import { Minus, Plus, ShoppingCart, Trash2 } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { FREE_SHIPPING_THRESHOLD } from '@ecommerce/lib/cart';
	import { formatPrice } from '@ecommerce/lib/format';
	import { cart } from '#lib/stores/cart.svelte';
	import { Button } from '@ecommerce/ui/button';
	import { Separator } from '@ecommerce/ui/separator';
	import { Skeleton } from '@ecommerce/ui/skeleton';

	const remainingForFreeShipping = $derived(Math.max(0, FREE_SHIPPING_THRESHOLD - cart.totalPrice));
</script>

<svelte:head>
	<title>{m.cart_metaTitle()}</title>
</svelte:head>

<h1 class="mb-6 text-2xl font-bold">
	{#if cart.hydrated && cart.items.length > 0}
		{m.cart_title({ count: String(cart.totalItems) })}
	{:else}
		{m.cart_metaTitle()}
	{/if}
</h1>

{#if !cart.hydrated}
	<div class="grid gap-3">
		<Skeleton class="h-24 w-full" />
		<Skeleton class="h-24 w-full" />
		<Skeleton class="h-24 w-full" />
	</div>
{:else if cart.items.length === 0}
	<div
		class="grid justify-items-center gap-4 rounded-xl border border-border bg-card p-12 text-center"
	>
		<ShoppingCart class="size-12 text-muted-foreground" />
		<div>
			<p class="text-lg font-semibold">{m.cart_emptyTitle()}</p>
			<p class="mt-1 text-sm text-muted-foreground">{m.cart_emptyHint()}</p>
		</div>
		<Button href={localizeHref('/products')}>{m.cart_continueShopping()}</Button>
	</div>
{:else}
	<div class="grid gap-6 lg:grid-cols-[1fr_360px]">
		<div class="grid content-start gap-3">
			<div class="flex justify-end">
				<Button variant="ghost" size="sm" onclick={() => cart.clear()}>
					<Trash2 class="size-4" />
					{m.cart_clearAll()}
				</Button>
			</div>

			{#each cart.items as item (item.product.productId)}
				<div class="flex gap-4 rounded-xl border border-border bg-card p-4">
					<a
						href={localizeHref(`/products/${item.product.productId}`)}
						class="block size-20 shrink-0 overflow-hidden rounded-lg bg-muted"
					>
						{#if item.product.imageUrl}
							<img
								src={item.product.imageUrl}
								alt={item.product.productTitle}
								class="h-full w-full object-cover"
							/>
						{:else}
							<div
								class="flex h-full w-full items-center justify-center text-3xl text-muted-foreground"
							>
								📦
							</div>
						{/if}
					</a>

					<div class="flex min-w-0 flex-1 flex-col gap-1">
						<a
							href={localizeHref(`/products/${item.product.productId}`)}
							class="line-clamp-2 text-sm font-medium hover:text-primary"
						>
							{item.product.productTitle}
						</a>
						<span class="font-semibold text-primary">{formatPrice(item.product.priceUnit)}</span>

						<div class="mt-auto flex items-center gap-3">
							<div class="flex items-center rounded-lg border border-border">
								<Button
									variant="ghost"
									size="icon-xs"
									aria-label={m.common_decrease()}
									onclick={() => cart.updateQuantity(item.product.productId, item.quantity - 1)}
								>
									<Minus class="size-3" />
								</Button>
								<span class="w-8 text-center text-sm">{item.quantity}</span>
								<Button
									variant="ghost"
									size="icon-xs"
									aria-label={m.common_increase()}
									onclick={() => cart.updateQuantity(item.product.productId, item.quantity + 1)}
								>
									<Plus class="size-3" />
								</Button>
							</div>
							<Button
								variant="ghost"
								size="sm"
								class="text-destructive"
								onclick={() => cart.remove(item.product.productId)}
							>
								<Trash2 class="size-4" />
							</Button>
						</div>
					</div>

					<div class="text-right text-sm font-semibold">
						{formatPrice(item.product.priceUnit * item.quantity)}
					</div>
				</div>
			{/each}
		</div>

		<aside class="h-fit rounded-xl border border-border bg-card p-5">
			<h2 class="mb-4 font-semibold">{m.cart_orderSummary()}</h2>

			<div class="grid gap-2 text-sm">
				<div class="flex justify-between">
					<span class="text-muted-foreground">
						{m.cart_subtotal({ count: String(cart.totalItems) })}
					</span>
					<span>{formatPrice(cart.totalPrice)}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-muted-foreground">{m.cart_shippingFee()}</span>
					<span>
						{#if cart.shippingFee === 0}
							<span class="text-green-600">{m.cart_freeShipping()}</span>
						{:else}
							{formatPrice(cart.shippingFee)}
						{/if}
					</span>
				</div>
			</div>

			{#if remainingForFreeShipping > 0}
				<p class="mt-3 rounded-lg bg-muted p-2 text-xs text-muted-foreground">
					{m.cart_shippingPromo({ amount: formatPrice(remainingForFreeShipping) })}
				</p>
			{/if}

			<Separator class="my-4" />

			<div class="flex justify-between text-base font-bold">
				<span>{m.cart_total()}</span>
				<span class="text-primary">{formatPrice(cart.grandTotal)}</span>
			</div>

			<Button href={localizeHref('/checkout')} size="lg" class="mt-4 w-full">
				{m.cart_checkout()}
			</Button>

			<p class="mt-3 text-center text-xs text-muted-foreground">{m.cart_securePayment()}</p>
		</aside>
	</div>
{/if}
