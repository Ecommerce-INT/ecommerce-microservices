<script lang="ts">
	import { Headphones, LogOut, MapPin, Package, Search, ShoppingCart, User } from '@lucide/svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import type { SessionUser } from '@ecommerce/lib/types';
	import { m } from '$lib/paraglide/messages.js';
	import { getLocale, locales, localizeHref } from '$lib/paraglide/runtime.js';
	import { cart } from '$lib/stores/cart.svelte';
	import { logout } from '$lib/remote/auth.remote';
	import { Badge } from '@ecommerce/ui/badge';
	import { Button } from '@ecommerce/ui/button';
	import * as DropdownMenu from '@ecommerce/ui/dropdown-menu';
	import { Input } from '@ecommerce/ui/input';

	let { user = null }: { user?: SessionUser | null } = $props();

	const locale = $derived(getLocale());
	const cartCount = $derived(cart.hydrated ? cart.totalItems : 0);
	const searchAction = $derived(localizeHref('/products'));
	let loggingOut = $state(false);

	const navCategories = [
		{ href: '/products?categoryId=1', label: m.category_c1 },
		{ href: '/products?categoryId=2', label: m.category_c2 },
		{ href: '/products?categoryId=3', label: m.category_c3 },
		{ href: '/products?categoryId=4', label: m.category_c4 },
		{ href: '/products?categoryId=5', label: m.category_c5 },
		{ href: '/products?categoryId=6', label: m.category_c6 },
		{ href: '/products?categoryId=7', label: m.category_c7 },
		{ href: '/products?categoryId=8', label: m.category_c8 },
		{ href: '/products?categoryId=9', label: m.category_c9 },
		{ href: '/products?categoryId=10', label: m.category_c10 }
	];

	async function handleLogout() {
		if (loggingOut) return;
		loggingOut = true;
		try {
			await logout();
			await invalidateAll();
			await goto(localizeHref('/'));
		} finally {
			loggingOut = false;
		}
	}
</script>

<header class="sticky top-0 z-40 shadow-sm">
	<div class="bg-foreground text-xs text-background">
		<div class="mx-auto flex h-8 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
			<div class="flex items-center gap-4">
				<span class="animate-brand-pulse hidden sm:inline">{m.header_topShipping()}</span>
				<a href={localizeHref('/orders')} class="flex items-center gap-1 hover:text-primary">
					<MapPin class="size-3" />
					{m.header_topTrackOrder()}
				</a>
			</div>
			<div class="flex items-center gap-4">
				<span class="hidden items-center gap-1 sm:flex">
					<Headphones class="size-3" />
					{m.header_hotline()}
				</span>
				<div class="flex items-center gap-1">
					{#each locales as l (l)}
						{#if l === locale}
							<span class="font-semibold text-primary">{l.toUpperCase()}</span>
						{:else}
							<a
								href={localizeHref(page.url.pathname, { locale: l })}
								data-sveltekit-reload
								class="hover:text-primary"
							>
								{l.toUpperCase()}
							</a>
						{/if}
					{/each}
				</div>
			</div>
		</div>
	</div>

	<div class="bg-card">
		<div class="mx-auto flex max-w-7xl items-center gap-4 px-4 py-3 sm:px-6 lg:px-8">
			<a href={localizeHref('/')} class="shrink-0 text-2xl font-black tracking-tight text-primary">
				EzBuy
			</a>

			<form action={searchAction} method="get" class="hidden flex-1 items-center gap-2 sm:flex">
				<Input
					name="search"
					placeholder={m.header_searchPlaceholder()}
					aria-label={m.header_searchPlaceholder()}
					class="h-9 flex-1"
				/>
				<Button type="submit" size="icon-sm" aria-label={m.header_searchPlaceholder()}>
					<Search class="size-4" />
				</Button>
			</form>

			<div class="ml-auto flex items-center gap-2">
				{#if user}
					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class="inline-flex h-8 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium hover:bg-muted"
						>
							<User class="size-4" />
							<span class="hidden max-w-24 truncate sm:inline">
								{user.fullName || user.username || m.header_myAccount()}
							</span>
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end" class="w-52">
							<DropdownMenu.Label>{user.email || user.username}</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.Item onclick={() => goto(localizeHref('/orders'))}>
								<Package class="size-4" />
								{m.header_myOrders()}
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => goto(localizeHref('/cart'))}>
								<ShoppingCart class="size-4" />
								{m.header_cart()}
							</DropdownMenu.Item>
							<DropdownMenu.Separator />
							<DropdownMenu.Item
								class="text-destructive"
								disabled={loggingOut}
								onclick={handleLogout}
							>
								<LogOut class="size-4" />
								{m.header_logout()}
							</DropdownMenu.Item>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{:else}
					<Button href={localizeHref('/login')} variant="ghost" size="sm" class="gap-1.5">
						<User class="size-4" />
						<span class="hidden sm:inline">{m.header_login()}</span>
					</Button>
				{/if}

				<Button href={localizeHref('/cart')} variant="outline" size="sm" class="relative gap-1.5">
					<ShoppingCart class="size-4" />
					<span class="hidden sm:inline">{m.header_cart()}</span>
					{#if cartCount > 0}
						<Badge class="absolute -top-2 -right-2 h-5 min-w-5 justify-center px-1">
							{cartCount}
						</Badge>
					{/if}
				</Button>
			</div>
		</div>
	</div>

	<nav class="hidden border-t border-border bg-card md:block">
		<div
			class="mx-auto flex max-w-7xl items-center gap-5 overflow-x-auto px-4 py-2 text-sm sm:px-6 lg:px-8"
		>
			{#each navCategories as item (item.href)}
				<a href={localizeHref(item.href)} class="whitespace-nowrap hover:text-primary">
					{item.label()}
				</a>
			{/each}
		</div>
	</nav>
</header>
