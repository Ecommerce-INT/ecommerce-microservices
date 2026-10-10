<script lang="ts">
	import { Pencil, Plus, Trash2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import type { Category } from '@ecommerce/lib/types';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import * as Dialog from '@ecommerce/ui/dialog';
	import { Input } from '@ecommerce/ui/input';
	import { Label } from '@ecommerce/ui/label';
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import { deleteCategory, getCategories, saveCategory } from './categories.remote';

	const currentPage = $derived(Math.max(0, Number(page.url.searchParams.get('page') ?? 0) || 0));
	const args = $derived({ page: currentPage, size: 12 });

	let dialogOpen = $state(false);
	let editing = $state<Category | null>(null);
	let deleteTarget = $state<Category | null>(null);
	let deleting = $state(false);

	const href = (target: number) =>
		localizeHref(target === 0 ? '/categories' : `/categories?page=${target}`);

	function openCreate() {
		editing = null;
		dialogOpen = true;
	}

	function openEdit(category: Category) {
		editing = category;
		dialogOpen = true;
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleting = true;
		try {
			await deleteCategory(deleteTarget.categoryId).updates(getCategories);
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
	<title>{m.meta_template({ page: m.categories_metaTitle() })}</title>
</svelte:head>

<div class="mb-4 flex flex-wrap items-center gap-3">
	<h1 class="text-xl font-bold">
		{#await getCategories(args)}
			{m.categories_title({ count: '…' })}
		{:then result}
			{m.categories_title({ count: String(result.totalElements) })}
		{/await}
	</h1>
	<Button size="sm" class="ml-auto" onclick={openCreate}>
		<Plus class="size-4" />
		{m.categories_addBtn()}
	</Button>
</div>

{#await getCategories(args)}
	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
		{#each Array(4).keys() as i (i)}
			<Skeleton class="h-40 rounded-xl" />
		{/each}
	</div>
{:then result}
	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
		{#each result.content as category (category.categoryId)}
			<article class="overflow-hidden rounded-xl border border-border bg-card">
				<div class="h-28 bg-muted">
					{#if category.imageUrl}
						<img
							src={category.imageUrl}
							alt={category.categoryTitle}
							class="h-full w-full object-cover"
						/>
					{/if}
				</div>
				<div class="flex items-start gap-2 p-4">
					<div class="min-w-0 flex-1">
						<p class="truncate font-medium">{category.categoryTitle}</p>
						<p class="text-xs text-muted-foreground">
							{m.categories_idLabel({ id: String(category.categoryId) })}
						</p>
					</div>
					<Button
						variant="ghost"
						size="icon-sm"
						aria-label={m.common_edit()}
						onclick={() => openEdit(category)}
					>
						<Pencil class="size-4" />
					</Button>
					<Button
						variant="ghost"
						size="icon-sm"
						class="text-destructive"
						aria-label={m.common_delete()}
						onclick={() => (deleteTarget = category)}
					>
						<Trash2 class="size-4" />
					</Button>
				</div>
			</article>
		{/each}
	</div>

	{#if result.content.length === 0}
		<p class="p-10 text-center text-sm text-muted-foreground">{m.common_noData()}</p>
	{/if}

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
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>
				{editing ? m.categories_modalEditTitle() : m.categories_modalAddTitle()}
			</Dialog.Title>
		</Dialog.Header>

		<form
			{...saveCategory.enhance(async (form) => {
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
				{...saveCategory.fields.categoryId.as('hidden', editing ? String(editing.categoryId) : '')}
			/>

			<div class="grid gap-4">
				<div class="grid gap-1.5">
					<Label for="categoryTitle">{m.categories_nameLabel()}</Label>
					<Input
						id="categoryTitle"
						{...saveCategory.fields.categoryTitle.as('text', editing?.categoryTitle ?? '')}
					/>
					{#each saveCategory.fields.categoryTitle.issues() ?? [] as issue, index (index)}
						<p class="text-xs text-destructive">{issue.message}</p>
					{/each}
				</div>
				<div class="grid gap-1.5">
					<Label for="imageUrl">{m.categories_imageLabel()}</Label>
					<Input
						id="imageUrl"
						{...saveCategory.fields.imageUrl.as('text', editing?.imageUrl ?? '')}
					/>
				</div>
			</div>

			<Dialog.Footer class="mt-5 gap-2">
				<Button type="button" variant="outline" onclick={() => (dialogOpen = false)}>
					{m.common_cancel()}
				</Button>
				<Button type="submit" disabled={!!saveCategory.pending}>
					{editing ? m.common_update() : m.common_save()}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root open={deleteTarget !== null} onOpenChange={(open) => !open && (deleteTarget = null)}>
	<Dialog.Content class="sm:max-w-sm">
		<Dialog.Header>
			<Dialog.Title>{m.categories_deleteConfirm()}</Dialog.Title>
		</Dialog.Header>
		<Dialog.Footer class="gap-2">
			<Button variant="outline" onclick={() => (deleteTarget = null)}>{m.common_cancel()}</Button>
			<Button variant="destructive" disabled={deleting} onclick={confirmDelete}>
				{m.common_delete()}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
