<script lang="ts">
	import { CircleCheckBig, CreditCard, Loader2, ShoppingCart, Wallet } from '@lucide/svelte';
	import { untrack } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { formatPrice } from '@ecommerce/lib/format';
	import { shippingSchema } from '@ecommerce/lib/schemas';
	import { m } from '$lib/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { cart } from '$lib/stores/cart.svelte';
	import { Button } from '@ecommerce/ui/button';
	import { Input } from '@ecommerce/ui/input';
	import { Label } from '@ecommerce/ui/label';
	import { Separator } from '@ecommerce/ui/separator';
	import { Textarea } from '@ecommerce/ui/textarea';
	import { placeOrder } from './checkout.remote';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let step = $state(1);
	let showErrors = $state(false);
	let placing = $state(false);
	let orderId = $state<number | null>(null);
	let paymentMethod = $state<'COD' | 'ONLINE'>('COD');
	let shipping = $state({
		fullName: untrack(() => data.user)?.fullName ?? '',
		phone: '',
		address: '',
		city: '',
		district: '',
		note: ''
	});

	const shippingValid = $derived(shippingSchema.safeParse(shipping).success);
	const steps = [
		{ number: 1, label: m.checkout_stepShipping },
		{ number: 2, label: m.checkout_stepPayment },
		{ number: 3, label: m.checkout_stepConfirm }
	];

	function nextFromShipping() {
		showErrors = true;
		if (shippingValid) step = 2;
	}

	async function submit() {
		if (placing) return;
		placing = true;
		try {
			const result = await placeOrder({
				shipping: { ...shipping },
				paymentMethod,
				items: cart.items.map((item) => ({
					productId: item.product.productId,
					quantity: item.quantity
				}))
			});
			orderId = result.orderId;
			cart.clear();
			step = 4;
		} catch {
			toast.error(m.checkout_orderFailed());
		} finally {
			placing = false;
		}
	}
</script>

<svelte:head>
	<title>{m.checkout_metaTitle()}</title>
</svelte:head>

<h1 class="mb-6 text-2xl font-bold">{m.checkout_title()}</h1>

{#if cart.hydrated && cart.items.length === 0 && step !== 4}
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
	<div class="mb-6 flex items-center gap-2">
		{#each steps as item (item.number)}
			<div class="flex items-center gap-2">
				<span
					class="flex size-7 items-center justify-center rounded-full text-sm font-semibold {step >=
					item.number
						? 'bg-primary text-primary-foreground'
						: 'bg-muted text-muted-foreground'}"
				>
					{item.number}
				</span>
				<span class="text-sm {step === item.number ? 'font-semibold' : 'text-muted-foreground'}">
					{item.label()}
				</span>
				{#if item.number < 3}
					<span class="mx-1 h-px w-8 bg-border"></span>
				{/if}
			</div>
		{/each}
	</div>

	<div class="grid gap-6 lg:grid-cols-[1fr_340px]">
		<div class="rounded-xl border border-border bg-card p-5">
			{#if step === 1}
				<h2 class="mb-4 font-semibold">{m.checkout_stepShipping()}</h2>
				<div class="grid gap-4 sm:grid-cols-2">
					<div class="grid gap-1.5">
						<Label for="fullName">{m.checkout_fullName()}</Label>
						<Input
							id="fullName"
							bind:value={shipping.fullName}
							aria-invalid={showErrors && !shippingValid && shipping.fullName.trim().length < 2}
						/>
					</div>
					<div class="grid gap-1.5">
						<Label for="phone">{m.checkout_phone()}</Label>
						<Input id="phone" bind:value={shipping.phone} inputmode="tel" />
					</div>
					<div class="grid gap-1.5 sm:col-span-2">
						<Label for="address">{m.checkout_address()}</Label>
						<Input
							id="address"
							bind:value={shipping.address}
							placeholder={m.checkout_addressPlaceholder()}
						/>
					</div>
					<div class="grid gap-1.5">
						<Label for="city">{m.checkout_city()}</Label>
						<Input id="city" bind:value={shipping.city} placeholder={m.checkout_selectCity()} />
					</div>
					<div class="grid gap-1.5">
						<Label for="district">{m.checkout_district()}</Label>
						<Input id="district" bind:value={shipping.district} />
					</div>
					<div class="grid gap-1.5 sm:col-span-2">
						<Label for="note">{m.checkout_note()}</Label>
						<Textarea
							id="note"
							bind:value={shipping.note}
							placeholder={m.checkout_notePlaceholder()}
						/>
					</div>
				</div>

				{#if showErrors && !shippingValid}
					<p class="mt-3 text-sm text-destructive">{m.checkout_validationErr()}</p>
				{/if}

				<div class="mt-5 flex justify-end">
					<Button onclick={nextFromShipping}>{m.checkout_continueToPayment()}</Button>
				</div>
			{:else if step === 2}
				<h2 class="mb-4 font-semibold">{m.checkout_paymentMethod()}</h2>
				<div class="grid gap-3">
					<button
						type="button"
						class="flex items-center gap-3 rounded-xl border p-4 text-left transition-colors {paymentMethod ===
						'COD'
							? 'border-primary bg-primary/5'
							: 'border-border hover:border-primary/50'}"
						onclick={() => (paymentMethod = 'COD')}
					>
						<Wallet class="size-5 shrink-0 text-primary" />
						<span>
							<span class="block text-sm font-medium">{m.checkout_cod()}</span>
							<span class="block text-xs text-muted-foreground">{m.checkout_codDesc()}</span>
						</span>
					</button>
					<button
						type="button"
						class="flex items-center gap-3 rounded-xl border p-4 text-left transition-colors {paymentMethod ===
						'ONLINE'
							? 'border-primary bg-primary/5'
							: 'border-border hover:border-primary/50'}"
						onclick={() => (paymentMethod = 'ONLINE')}
					>
						<CreditCard class="size-5 shrink-0 text-primary" />
						<span>
							<span class="block text-sm font-medium">{m.checkout_bankTransfer()}</span>
							<span class="block text-xs text-muted-foreground"
								>{m.checkout_bankTransferDesc()}</span
							>
						</span>
					</button>
				</div>

				<div class="mt-5 flex justify-between">
					<Button variant="outline" onclick={() => (step = 1)}>{m.checkout_back()}</Button>
					<Button onclick={() => (step = 3)}>{m.checkout_continueToPayment()}</Button>
				</div>
			{:else if step === 3}
				<h2 class="mb-4 font-semibold">{m.checkout_stepConfirm()}</h2>
				<dl class="grid gap-2 text-sm">
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">{m.checkout_fullName()}</dt>
						<dd class="text-right font-medium">{shipping.fullName}</dd>
					</div>
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">{m.checkout_phone()}</dt>
						<dd class="text-right font-medium">{shipping.phone}</dd>
					</div>
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">{m.checkout_address()}</dt>
						<dd class="text-right font-medium">
							{shipping.address}, {shipping.district}, {shipping.city}
						</dd>
					</div>
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">{m.checkout_paymentMethod()}</dt>
						<dd class="text-right font-medium">
							{paymentMethod === 'COD' ? m.checkout_cod() : m.checkout_bankTransfer()}
						</dd>
					</div>
				</dl>

				<div class="mt-5 flex justify-between">
					<Button variant="outline" onclick={() => (step = 2)}>{m.checkout_back()}</Button>
					<Button onclick={submit} disabled={placing}>
						{#if placing}
							<Loader2 class="size-4 animate-spin" />
						{/if}
						{m.checkout_placeOrder()}
					</Button>
				</div>
			{:else}
				<div class="grid justify-items-center gap-3 py-6 text-center">
					<CircleCheckBig class="size-12 text-green-600" />
					<h2 class="text-lg font-bold">{m.checkout_successTitle()}</h2>
					<p class="text-sm">
						{m.checkout_orderCode()} <span class="font-semibold">#{orderId}</span>
					</p>
					<p class="text-sm text-muted-foreground">{m.checkout_successMsg()}</p>
					<div class="mt-2 flex flex-wrap justify-center gap-2">
						<Button href={localizeHref('/orders')}>{m.checkout_viewOrders()}</Button>
						<Button href={localizeHref('/products')} variant="outline">
							{m.checkout_continueShopping()}
						</Button>
					</div>
				</div>
			{/if}
		</div>

		{#if step !== 4}
			<aside class="h-fit rounded-xl border border-border bg-card p-5">
				<h2 class="mb-4 font-semibold">{m.checkout_yourOrder()}</h2>
				<div class="grid gap-2 text-sm">
					{#each cart.items as item (item.product.productId)}
						<div class="flex justify-between gap-3">
							<span class="line-clamp-1 text-muted-foreground">
								{item.product.productTitle} × {item.quantity}
							</span>
							<span>{formatPrice(item.product.priceUnit * item.quantity)}</span>
						</div>
					{/each}
				</div>

				<Separator class="my-4" />

				<div class="grid gap-2 text-sm">
					<div class="flex justify-between">
						<span class="text-muted-foreground">{m.checkout_subtotal()}</span>
						<span>{formatPrice(cart.totalPrice)}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-muted-foreground">{m.checkout_shipping()}</span>
						<span>
							{#if cart.shippingFee === 0}
								<span class="text-green-600">{m.checkout_freeShipping()}</span>
							{:else}
								{formatPrice(cart.shippingFee)}
							{/if}
						</span>
					</div>
				</div>

				<Separator class="my-4" />

				<div class="flex justify-between font-bold">
					<span>{m.checkout_total()}</span>
					<span class="text-primary">{formatPrice(cart.grandTotal)}</span>
				</div>
			</aside>
		{/if}
	</div>
{/if}
