<script lang="ts">
	import { PackageOpen, Pencil, Plus, Search, Trash2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import { formatPrice } from '@ecommerce/lib/format';
	import type { Product } from '@ecommerce/lib/types';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import * as Dialog from '@ecommerce/ui/dialog';
	import { EmptyState } from '@ecommerce/ui/empty-state';
	import { Input } from '@ecommerce/ui/input';
	import { Label } from '@ecommerce/ui/label';
	import * as Table from '@ecommerce/ui/table';
	import { TableSkeleton } from '@ecommerce/ui/table-skeleton';
	import { deleteProduct, getCategoryOptions, getProducts, saveProduct } from './products.remote';

	const currentPage = $derived(Math.max(0, Number(page.url.searchParams.get('page') ?? 0) || 0));
	const search = $derived(page.url.searchParams.get('search')?.trim() || undefined);
	const args = $derived({ page: currentPage, size: 10, search });

	let dialogOpen = $state(false);
	let editing = $state<Product | null>(null);
	let deleteTarget = $state<Product | null>(null);
	let deleting = $state(false);

	const href = (target: number) =>
		localizeHref(target === 0 ? '/products' : `/products?page=${target}`);

	function openCreate() {
		editing = null;
		dialogOpen = true;
	}

	function openEdit(product: Product) {
		editing = product;
		dialogOpen = true;
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleting = true;
		try {
			await deleteProduct(deleteTarget.productId).updates(getProducts);
			toast.success(m.common_delete());
			deleteTarget = null;
		} catch {
			toast.error(m.common_noData());
		} finally {
			deleting = false;
		}
	}
</script>

<svelte:head>
	<title>{m.meta_template({ page: m.products_metaTitle() })}</title>
</svelte:head>

<div class="mb-4 flex flex-wrap items-center gap-3">
	<h1 class="text-xl font-bold">{m.products_title()}</h1>

	<form method="get" action={localizeHref('/products')} class="ml-auto flex items-center gap-2">
		<Input
			name="search"
			value={search ?? ''}
			placeholder={m.products_searchPlaceholder()}
			class="h-9 w-52"
		/>
		<Button type="submit" variant="outline" size="icon-sm" aria-label={m.common_search()}>
			<Search class="size-4" />
		</Button>
	</form>

	<Button size="sm" onclick={openCreate}>
		<Plus class="size-4" />
		{m.products_addBtn()}
	</Button>
</div>

{#await getProducts(args)}
	<TableSkeleton rows={6} />
{:then result}
	<p class="mb-2 text-sm text-muted-foreground">
		{m.products_count({ count: String(result.totalElements) })}
	</p>

	<div class="overflow-hidden rounded-xl border border-border bg-card">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>{m.products_colProduct()}</Table.Head>
					<Table.Head>{m.products_colSku()}</Table.Head>
					<Table.Head>{m.products_colCategory()}</Table.Head>
					<Table.Head class="text-right">{m.products_colPrice()}</Table.Head>
					<Table.Head class="text-center">{m.products_colStock()}</Table.Head>
					<Table.Head class="text-right">{m.products_colActions()}</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each result.content as product (product.productId)}
					<Table.Row>
						<Table.Cell>
							<div class="flex items-center gap-3">
								<div class="size-9 shrink-0 overflow-hidden rounded-lg bg-muted">
									{#if product.imageUrl}
										<img
											src={product.imageUrl}
											alt={product.productTitle}
											class="h-full w-full object-cover"
										/>
									{/if}
								</div>
								<span class="line-clamp-1 font-medium">{product.productTitle}</span>
							</div>
						</Table.Cell>
						<Table.Cell class="text-muted-foreground">{product.sku ?? '—'}</Table.Cell>
						<Table.Cell class="text-muted-foreground">
							{product.category?.categoryTitle ?? m.products_uncategorized()}
						</Table.Cell>
						<Table.Cell class="text-right font-medium">{formatPrice(product.priceUnit)}</Table.Cell>
						<Table.Cell class="text-center">
							{#if product.quantity > 0}
								{product.quantity}
							{:else}
								<Badge variant="destructive">{m.products_outOfStock()}</Badge>
							{/if}
						</Table.Cell>
						<Table.Cell>
							<div class="flex justify-end gap-1">
								<Button
									variant="ghost"
									size="icon-sm"
									aria-label={m.common_edit()}
									onclick={() => openEdit(product)}
								>
									<Pencil class="size-4" />
								</Button>
								<Button
									variant="ghost"
									size="icon-sm"
									class="text-destructive"
									aria-label={m.common_delete()}
									onclick={() => (deleteTarget = product)}
								>
									<Trash2 class="size-4" />
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>

		{#if result.content.length === 0}
			<EmptyState
				class="rounded-none border-0"
				icon={PackageOpen}
				title={m.products_noProducts()}
			/>
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

<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>
				{editing ? m.products_modalEditTitle() : m.products_modalAddTitle()}
			</Dialog.Title>
		</Dialog.Header>

		<form
			{...saveProduct.enhance(async (form) => {
				try {
					if (await form.submit()) {
						dialogOpen = false;
						editing = null;
						form.element.reset();
						toast.success(m.common_save());
					}
				} catch {
					toast.error(m.common_noData());
				}
			})}
		>
			<input
				{...saveProduct.fields.productId.as('hidden', editing ? String(editing.productId) : '')}
			/>

			<div class="grid gap-4">
				<div class="grid gap-1.5">
					<Label for="productTitle">{m.products_nameLabel()}</Label>
					<Input
						id="productTitle"
						{...saveProduct.fields.productTitle.as('text', editing?.productTitle ?? '')}
					/>
					{#each saveProduct.fields.productTitle.issues() ?? [] as issue, index (index)}
						<p class="text-xs text-destructive">{issue.message}</p>
					{/each}
				</div>

				<div class="grid gap-4 sm:grid-cols-2">
					<div class="grid gap-1.5">
						<Label for="sku">{m.products_skuLabel()}</Label>
						<Input id="sku" {...saveProduct.fields.sku.as('text', editing?.sku ?? '')} />
						{#each saveProduct.fields.sku.issues() ?? [] as issue, index (index)}
							<p class="text-xs text-destructive">{issue.message}</p>
						{/each}
					</div>
					<div class="grid gap-1.5">
						<Label for="categoryId">{m.products_categoryLabel()}</Label>
						<select
							id="categoryId"
							name="categoryId"
							class="h-9 rounded-lg border border-border bg-background px-2 text-sm"
						>
							<option value="">{m.products_selectCategory()}</option>
							{#await getCategoryOptions()}
								<option value="">{m.common_loading()}</option>
							{:then categories}
								{#each categories as category (category.categoryId)}
									<option
										value={category.categoryId}
										selected={editing?.category?.categoryId === category.categoryId}
									>
										{category.categoryTitle}
									</option>
								{/each}
							{/await}
						</select>
					</div>
				</div>

				<div class="grid gap-4 sm:grid-cols-2">
					<div class="grid gap-1.5">
						<Label for="priceUnit">{m.products_priceLabel()}</Label>
						<Input
							id="priceUnit"
							{...saveProduct.fields.priceUnit.as('text', editing ? String(editing.priceUnit) : '')}
							inputmode="numeric"
						/>
						{#each saveProduct.fields.priceUnit.issues() ?? [] as issue, index (index)}
							<p class="text-xs text-destructive">{issue.message}</p>
						{/each}
					</div>
					<div class="grid gap-1.5">
						<Label for="quantity">{m.products_quantityLabel()}</Label>
						<Input
							id="quantity"
							{...saveProduct.fields.quantity.as('text', editing ? String(editing.quantity) : '')}
							inputmode="numeric"
						/>
						{#each saveProduct.fields.quantity.issues() ?? [] as issue, index (index)}
							<p class="text-xs text-destructive">{issue.message}</p>
						{/each}
					</div>
				</div>

				<div class="grid gap-1.5">
					<Label for="imageUrl">{m.products_imageLabel()}</Label>
					<Input
						id="imageUrl"
						{...saveProduct.fields.imageUrl.as('text', editing?.imageUrl ?? '')}
					/>
				</div>
			</div>

			<Dialog.Footer class="mt-5 gap-2">
				<Button type="button" variant="outline" onclick={() => (dialogOpen = false)}>
					{m.common_cancel()}
				</Button>
				<Button type="submit" disabled={!!saveProduct.pending}>
					{editing ? m.common_update() : m.common_save()}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root open={deleteTarget !== null} onOpenChange={(open) => !open && (deleteTarget = null)}>
	<Dialog.Content class="sm:max-w-sm">
		<Dialog.Header>
			<Dialog.Title>{m.products_deleteConfirm()}</Dialog.Title>
		</Dialog.Header>
		<Dialog.Footer class="gap-2">
			<Button variant="outline" onclick={() => (deleteTarget = null)}>{m.common_cancel()}</Button>
			<Button variant="destructive" disabled={deleting} onclick={confirmDelete}>
				{m.common_delete()}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
