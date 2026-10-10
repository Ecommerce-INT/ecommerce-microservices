<script lang="ts">
	import { Search, SlidersHorizontal } from '@lucide/svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Product } from '@ecommerce/lib/types';
	import { m } from '$lib/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import { Input } from '@ecommerce/ui/input';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import ProductCard from '$lib/components/product-card.svelte';
	import { getCategories, getProducts } from './catalog.remote';

	const SORTS = [
		{ value: 'productId,desc', label: m.products_sortNewest },
		{ value: 'priceUnit,asc', label: m.products_sortPriceLow },
		{ value: 'priceUnit,desc', label: m.products_sortPriceHigh }
	] as const;
	type SortValue = (typeof SORTS)[number]['value'];

	const PRICE_RANGES = [
		{ min: 0, max: 100_000, label: m.products_price1 },
		{ min: 100_000, max: 500_000, label: m.products_price2 },
		{ min: 500_000, max: 1_000_000, label: m.products_price3 },
		{ min: 1_000_000, max: 5_000_000, label: m.products_price4 },
		{ min: 5_000_000, max: Number.POSITIVE_INFINITY, label: m.products_price5 }
	];

	const searchParams = $derived(page.url.searchParams);
	const currentPage = $derived(Math.max(0, Number(searchParams.get('page') ?? 0) || 0));
	const search = $derived(searchParams.get('search')?.trim() || undefined);
	const categoryId = $derived(
		searchParams.get('categoryId') ? Number(searchParams.get('categoryId')) : undefined
	);
	const sort = $derived(
		SORTS.some((option) => option.value === searchParams.get('sort'))
			? (searchParams.get('sort') as SortValue)
			: undefined
	);
	const args = $derived({ page: currentPage, size: 12, search, categoryId, sort });

	let priceRange = $state<number | null>(null);

	function urlFor(overrides: Record<string, string | number | undefined>) {
		const params = new SvelteURLSearchParams(page.url.searchParams.toString());
		for (const [key, value] of Object.entries(overrides)) {
			if (value === undefined || value === '') params.delete(key);
			else params.set(key, String(value));
		}
		const qs = params.toString();
		return localizeHref(`/products${qs ? `?${qs}` : ''}`);
	}

	const searchAction = $derived(urlFor({ search: search ?? undefined, page: undefined }));

	function visibleProducts(products: Product[]) {
		const range = priceRange === null ? null : PRICE_RANGES[priceRange];
		if (!range) return products;
		return products.filter(
			(product) => product.priceUnit >= range.min && product.priceUnit < range.max
		);
	}

	function pageWindow(current: number, total: number) {
		const start = Math.max(0, Math.min(current - 2, total - 5));
		const end = Math.min(total, start + 5);
		return Array.from({ length: Math.max(0, end - start) }, (_, index) => start + index);
	}
</script>

<svelte:head>
	<title>{m.products_metaTitle()}</title>
</svelte:head>

<div class="grid gap-6 lg:grid-cols-[240px_1fr]">
	<aside class="grid content-start gap-6">
		<div class="rounded-xl border border-border bg-card p-4">
			<h2 class="mb-3 flex items-center gap-2 text-sm font-semibold">
				<SlidersHorizontal class="size-4" />
				{m.products_filtersTitle()}
			</h2>

			<p class="mb-2 text-xs font-medium text-muted-foreground uppercase">
				{m.products_categoriesTitle()}
			</p>
			<div class="grid gap-1">
				<a
					href={urlFor({ categoryId: undefined, page: undefined })}
					class="rounded px-2 py-1 text-sm hover:text-primary"
					class:bg-muted={categoryId === undefined}
					class:font-medium={categoryId === undefined}
				>
					{m.products_allCategories()}
				</a>
				{#await getCategories()}
					<div class="grid gap-1">
						<Skeleton class="h-6 w-full" />
						<Skeleton class="h-6 w-3/4" />
					</div>
				{:then categories}
					{#each categories as category (category.categoryId)}
						<a
							href={urlFor({ categoryId: category.categoryId, page: undefined })}
							class="rounded px-2 py-1 text-sm hover:text-primary"
							class:bg-muted={categoryId === category.categoryId}
							class:font-medium={categoryId === category.categoryId}
						>
							{category.categoryTitle}
						</a>
					{/each}
				{:catch}
					<p class="text-sm text-muted-foreground">—</p>
				{/await}
			</div>

			<p class="mt-4 mb-2 text-xs font-medium text-muted-foreground uppercase">
				{m.products_priceTitle()}
			</p>
			<div class="grid gap-1">
				{#each PRICE_RANGES as range, index (range.label)}
					<button
						type="button"
						class="rounded px-2 py-1 text-left text-sm hover:text-primary"
						class:bg-muted={priceRange === index}
						class:font-medium={priceRange === index}
						onclick={() => (priceRange = priceRange === index ? null : index)}
					>
						{range.label()}
					</button>
				{/each}
			</div>
		</div>
	</aside>

	<div>
		<div class="mb-4 flex flex-wrap items-center gap-3">
			<h1 class="text-xl font-bold">
				{#if search}
					{m.products_searchResults({ query: search })}
				{:else}
					{m.products_title()}
				{/if}
			</h1>

			<div class="ml-auto flex items-center gap-2">
				<form action={searchAction} method="get" class="flex items-center gap-2">
					<Input
						name="search"
						value={search ?? ''}
						placeholder={m.header_searchPlaceholder()}
						class="h-9 w-44"
					/>
					<Button type="submit" variant="outline" size="icon-sm" aria-label={m.products_filter()}>
						<Search class="size-4" />
					</Button>
				</form>

				<label class="text-sm text-muted-foreground" for="sort">
					{m.products_filter()}
				</label>
				<select
					id="sort"
					class="h-9 rounded-lg border border-border bg-card px-2 text-sm"
					value={sort ?? SORTS[0].value}
					onchange={(event) => goto(urlFor({ sort: event.currentTarget.value, page: undefined }))}
				>
					{#each SORTS as option (option.value)}
						<option value={option.value}>{option.label()}</option>
					{/each}
				</select>
			</div>
		</div>

		{#await getProducts(args)}
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-4">
				{#each Array(8).keys() as i (i)}
					<div class="grid gap-2">
						<Skeleton class="aspect-square w-full rounded-xl" />
						<Skeleton class="h-4 w-3/4" />
						<Skeleton class="h-4 w-1/2" />
					</div>
				{/each}
			</div>
		{:then result}
			{@const products = visibleProducts(result.content ?? [])}
			<p class="mb-3 text-sm text-muted-foreground">
				{m.products_count({ count: String(result.totalElements ?? 0) })}
			</p>

			{#if products.length === 0}
				<div
					class="grid justify-items-center gap-3 rounded-xl border border-border bg-card p-12 text-center"
				>
					<p class="font-semibold">{m.products_notFound()}</p>
					<p class="text-sm text-muted-foreground">{m.products_notFoundHint()}</p>
					<Button href={localizeHref('/products')} variant="outline" size="sm">
						{m.products_clearFilter()}
					</Button>
				</div>
			{:else}
				<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-4">
					{#each products as product (product.productId)}
						<ProductCard {product} />
					{/each}
				</div>

				{#if (result.totalPages ?? 0) > 1}
					<div class="mt-6 flex flex-wrap items-center justify-center gap-2">
						<Button
							variant="outline"
							size="sm"
							href={urlFor({ page: currentPage - 1 })}
							disabled={currentPage <= 0}
						>
							{m.common_prev()}
						</Button>

						{#each pageWindow(currentPage, result.totalPages) as pageNumber (pageNumber)}
							<Button
								variant={pageNumber === currentPage ? 'default' : 'outline'}
								size="icon-sm"
								href={urlFor({ page: pageNumber })}
								aria-current={pageNumber === currentPage ? 'page' : undefined}
							>
								{pageNumber + 1}
							</Button>
						{/each}

						<Button
							variant="outline"
							size="sm"
							href={urlFor({ page: currentPage + 1 })}
							disabled={currentPage >= (result.totalPages ?? 1) - 1}
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
				<p class="font-semibold">{m.productgrid_loadError()}</p>
				<Button href={localizeHref('/products')} variant="outline" size="sm">
					{m.common_retry()}
				</Button>
			</div>
		{/await}
	</div>
</div>
