<script lang="ts">
	import { onMount } from 'svelte';
	import { ChevronRight, Flame, RefreshCw, ShieldCheck, Truck } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Badge } from '@ecommerce/ui/badge';
	import ProductCard from '#lib/components/product-card.svelte';
	import {
		brands,
		categoryIcons,
		flashSaleProducts,
		heroSideBanners,
		heroSidebar,
		heroSlides,
		hotProducts,
		laptopBrands,
		laptopProducts,
		nextFlashSaleDeadline,
		phoneBrands,
		phoneProducts,
		quickDeals
	} from '#lib/data/home-mock';

	const slideText = [
		{ title: m.home_slide1Title, sub: m.home_slide1Sub },
		{ title: m.home_slide2Title, sub: m.home_slide2Sub },
		{ title: m.home_slide3Title, sub: m.home_slide3Sub }
	];
	const sideText = [
		{ title: m.home_side1Title, sub: m.home_side1Sub },
		{ title: m.home_side2Title, sub: m.home_side2Sub }
	];
	const dealText = [m.home_deal1, m.home_deal2, m.home_deal3, m.home_deal4, m.home_deal5];
	const categoryText = [
		m.category_c1,
		m.category_c2,
		m.category_c3,
		m.category_c4,
		m.category_c5,
		m.category_c6,
		m.category_c7,
		m.category_c8,
		m.category_c9,
		m.category_c10
	];
	const sidebarText = [
		m.category_c1,
		m.category_c2,
		m.category_c4,
		m.category_c5,
		m.category_c7,
		m.category_c9,
		m.category_c8,
		m.category_c10,
		m.category_c9,
		m.category_c1
	];
	const usp = [
		{ icon: ShieldCheck, title: m.home_uspGenuine, desc: m.home_uspGenuineDesc },
		{ icon: Truck, title: m.home_uspShipping, desc: m.home_uspShippingDesc },
		{ icon: RefreshCw, title: m.home_uspReturn, desc: m.home_uspReturnDesc },
		{ icon: Flame, title: m.home_uspSupport, desc: m.home_uspSupportDesc }
	];
	const productSections = [
		{
			href: '/products?categoryId=1',
			gradient: 'from-primary-600 to-primary-800',
			title: m.home_phoneBannerTitle,
			sub: m.home_phoneBannerSub,
			products: phoneProducts,
			brands: phoneBrands
		},
		{
			href: '/products?categoryId=2',
			gradient: 'from-slate-700 to-slate-900',
			title: m.home_laptopBannerTitle,
			sub: m.home_laptopBannerSub,
			products: laptopProducts,
			brands: laptopBrands
		}
	];

	let countdown = $state('');
	onMount(() => {
		const deadline = nextFlashSaleDeadline();
		const tick = () => {
			const ms = Math.max(0, deadline - Date.now());
			const hours = Math.floor(ms / 3_600_000);
			const minutes = Math.floor((ms % 3_600_000) / 60_000);
			const seconds = Math.floor((ms % 60_000) / 1000);
			countdown = [hours, minutes, seconds].map((n) => String(n).padStart(2, '0')).join(':');
		};
		tick();
		const id = setInterval(tick, 1000);
		return () => clearInterval(id);
	});
</script>

<svelte:head>
	<title>{m.home_metaTitle()}</title>
	<meta name="description" content={m.home_metaDesc()} />
</svelte:head>

<section class="grid gap-4 lg:grid-cols-[240px_1fr_260px]">
	<aside class="hidden overflow-hidden rounded-xl border border-border bg-card lg:block">
		<ul class="divide-y divide-border">
			{#each heroSidebar as item, index (item.id)}
				<li>
					<a
						href={localizeHref(item.href)}
						class="flex items-center gap-2 px-4 py-2 text-sm hover:text-primary"
					>
						<span>{item.emoji}</span>
						<span class="truncate">{sidebarText[index]()}</span>
					</a>
				</li>
			{/each}
		</ul>
	</aside>

	<div class="grid gap-4 sm:grid-cols-[1fr_240px]">
		<div class="grid gap-3">
			{#each heroSlides as slide, index (slide.id)}
				<a
					href={localizeHref(slide.href)}
					class="relative flex min-h-36 flex-col justify-center gap-1 rounded-xl bg-linear-to-br p-6 text-white {slide.gradient}"
				>
					<Badge class="w-fit bg-white/20 text-white">{m.home_heroBadge()}</Badge>
					<h2 class="text-xl font-bold sm:text-2xl">{slideText[index].title()}</h2>
					<p class="text-sm text-white/90">{slideText[index].sub()}</p>
					<span class="mt-1 flex items-center gap-1 text-sm font-semibold">
						{m.home_heroCta()}
						<ChevronRight class="size-4" />
					</span>
				</a>
			{/each}
		</div>

		<div class="hidden gap-4 sm:grid">
			{#each heroSideBanners as banner, index (banner.id)}
				<a
					href={localizeHref(banner.href)}
					class="flex flex-col justify-center rounded-xl bg-linear-to-br p-4 text-white {banner.gradient}"
				>
					<span class="font-semibold">{sideText[index].title()}</span>
					<span class="text-xs text-white/90">{sideText[index].sub()}</span>
				</a>
			{/each}
		</div>
	</div>
</section>

<section class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-5">
	{#each quickDeals as deal, index (deal.id)}
		<a
			href={localizeHref(deal.href)}
			class="flex items-center gap-3 rounded-lg border border-border bg-card px-3 py-2.5 text-sm font-medium transition-colors hover:border-primary"
		>
			<span class="text-lg">{deal.emoji}</span>
			{dealText[index]()}
		</a>
	{/each}
</section>

<section class="mt-8">
	<div class="grid grid-cols-5 gap-2 sm:grid-cols-10">
		{#each categoryIcons as icon, index (icon.id)}
			<a
				href={localizeHref(icon.href)}
				class="flex flex-col items-center gap-1.5 rounded-lg border border-border bg-card px-2 py-3 text-center text-xs transition-colors hover:border-primary"
			>
				<span class="text-xl">{icon.emoji}</span>
				<span class="line-clamp-2">{categoryText[index]()}</span>
			</a>
		{/each}
	</div>
</section>

<section class="mt-8 rounded-xl border border-border bg-card p-4">
	<div class="mb-4 flex flex-wrap items-center gap-3">
		<h2 class="flex items-center gap-2 text-lg font-bold">
			<Flame class="size-5 text-primary" />
			{m.home_flashSaleTitle()}
		</h2>
		<div class="flex items-center gap-2 text-sm">
			<span class="text-muted-foreground">{m.home_flashSaleEndsIn()}</span>
			<span class="rounded bg-foreground px-2 py-0.5 font-mono font-semibold text-background">
				{countdown || '00:00:00'}
			</span>
		</div>
		<a href={localizeHref('/products?sort=priceUnit,asc')} class="ml-auto text-sm text-primary">
			{m.common_viewAll()}
		</a>
	</div>

	<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
		{#each flashSaleProducts as product (product.productId)}
			<div class="flex flex-col gap-2">
				<ProductCard {product} />
				{#if product.stockTotal}
					{@const soldPercent = Math.min(
						100,
						Math.round(((product.sold ?? 0) / product.stockTotal) * 100)
					)}
					<div class="h-1.5 w-full overflow-hidden rounded-full bg-muted">
						<div class="h-full bg-primary" style="width: {soldPercent}%"></div>
					</div>
					<p class="text-center text-xs text-muted-foreground">
						{m.home_sold({ count: String(product.sold ?? 0) })}
					</p>
				{/if}
			</div>
		{/each}
	</div>
</section>

<section class="mt-8">
	<div class="mb-4 flex items-center justify-between">
		<h2 class="text-lg font-bold">{m.home_hotProducts()}</h2>
		<a href={localizeHref('/products')} class="text-sm text-primary">{m.common_viewAll()}</a>
	</div>
	<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
		{#each hotProducts as product (product.productId)}
			<ProductCard {product} />
		{/each}
	</div>
</section>

{#each productSections as section (section.href)}
	<section class="mt-8 grid gap-4 lg:grid-cols-[240px_1fr]">
		<a
			href={localizeHref(section.href)}
			class="flex flex-col justify-center gap-2 rounded-xl bg-linear-to-br p-5 text-white {section.gradient}"
		>
			<h2 class="text-xl font-bold">{section.title()}</h2>
			<p class="text-sm text-white/85">{section.sub()}</p>
			<span class="mt-1 flex items-center gap-1 text-sm font-semibold">
				{m.home_viewAll()}
				<ChevronRight class="size-4" />
			</span>
		</a>

		<div class="overflow-hidden">
			<div class="no-scrollbar mb-3 flex gap-2 overflow-x-auto pb-1">
				{#each section.brands as brand (brand.id)}
					<a
						href={localizeHref(brand.href)}
						class="rounded-full border border-border bg-card px-3 py-1 text-xs whitespace-nowrap hover:border-primary"
					>
						{brand.label}
					</a>
				{/each}
			</div>
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-5">
				{#each section.products as product (product.productId)}
					<ProductCard {product} />
				{/each}
			</div>
		</div>
	</section>
{/each}

<section class="mt-8">
	<h2 class="mb-4 text-lg font-bold">{m.home_featuredBrands()}</h2>
	<div class="grid grid-cols-4 gap-3 sm:grid-cols-8">
		{#each brands as brand (brand.id)}
			<a
				href={localizeHref(brand.href)}
				class="flex h-14 items-center justify-center rounded-lg border border-border bg-card text-sm font-semibold hover:border-primary"
			>
				{brand.name}
			</a>
		{/each}
	</div>
</section>

<section class="mt-8 grid grid-cols-2 gap-3 lg:grid-cols-4">
	{#each usp as item (item.title())}
		{@const Icon = item.icon}
		<div class="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
			<Icon class="size-6 shrink-0 text-primary" />
			<div>
				<p class="text-sm font-semibold">{item.title()}</p>
				<p class="text-xs text-muted-foreground">{item.desc()}</p>
			</div>
		</div>
	{/each}
</section>
