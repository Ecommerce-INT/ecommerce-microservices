<script lang="ts">
	import {
		Boxes,
		LayoutDashboard,
		LogOut,
		Menu,
		Package,
		Settings,
		ShoppingBag,
		Tags,
		Truck,
		Users
	} from '@lucide/svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { m } from '#lib/paraglide/messages.js';
	import { getLocale, locales, localizeHref } from '#lib/paraglide/runtime.js';
	import { Button } from '@ecommerce/ui/button';
	import * as Sheet from '@ecommerce/ui/sheet';
	import { Toaster } from '@ecommerce/ui/sonner';
	import { logout } from '#lib/remote/auth.remote';
	import './layout.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	const barePaths = $derived(
		['/login', '/auth/callback', '/access-denied'].map((route) => localizeHref(route))
	);
	const bare = $derived(
		barePaths.some(
			(route) => page.url.pathname === route || page.url.pathname.startsWith(`${route}/`)
		)
	);
	const locale = $derived(getLocale());
	let loggingOut = $state(false);
	let drawerOpen = $state(false);

	const navItems = [
		{ href: '/dashboard', label: m.sidebar_dashboard, icon: LayoutDashboard },
		{ href: '/products', label: m.sidebar_products, icon: Package },
		{ href: '/orders', label: m.sidebar_orders, icon: ShoppingBag },
		{ href: '/categories', label: m.sidebar_categories, icon: Tags },
		{ href: '/users', label: m.sidebar_users, icon: Users },
		{ href: '/inventory', label: m.sidebar_inventory, icon: Boxes },
		{ href: '/shipping', label: m.sidebar_shipping, icon: Truck },
		{ href: '/settings', label: m.sidebar_settings, icon: Settings }
	];

	function isActive(href: string) {
		const localized = localizeHref(href);
		return page.url.pathname === localized || page.url.pathname.startsWith(`${localized}/`);
	}

	async function handleLogout() {
		if (loggingOut) return;
		loggingOut = true;
		try {
			await logout();
			await invalidateAll();
			await goto(localizeHref('/login'));
		} finally {
			loggingOut = false;
		}
	}
</script>

<svelte:head>
	<title>{m.meta_title()}</title>
	<meta name="description" content={m.meta_description()} />
</svelte:head>

{#if bare}
	{@render children()}
{:else}
	<div class="flex min-h-screen">
		<aside class="hidden w-60 shrink-0 border-r border-border bg-card lg:flex lg:flex-col">
			<div class="border-b border-border px-5 py-4">
				<p class="text-xl font-black text-primary">{m.sidebar_brand()}</p>
				<p class="text-xs text-muted-foreground">{m.sidebar_brandSub()}</p>
			</div>
			<nav class="flex-1 p-2">
				{#each navItems as item (item.href)}
					{@const active = isActive(item.href)}
					<a
						href={localizeHref(item.href)}
						class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors {active
							? 'bg-muted font-semibold text-primary'
							: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
					>
						<item.icon class="size-4" />
						{item.label()}
					</a>
				{/each}
			</nav>
			<div class="border-t border-border p-3">
				<Button
					variant="ghost"
					size="sm"
					class="w-full justify-start"
					disabled={loggingOut}
					onclick={handleLogout}
				>
					<LogOut class="size-4" />
					{m.sidebar_logout()}
				</Button>
			</div>
		</aside>

		<div class="flex min-w-0 flex-1 flex-col">
			<header class="flex h-14 items-center gap-3 border-b border-border bg-card px-4 lg:px-6">
				<Button
					variant="ghost"
					size="icon-sm"
					class="lg:hidden"
					aria-label={m.sidebar_openMenu()}
					onclick={() => (drawerOpen = true)}
				>
					<Menu class="size-5" />
				</Button>
				<p class="text-lg font-black text-primary lg:hidden">{m.sidebar_brand()}</p>
				<div class="ml-auto flex items-center gap-2">
					<div class="flex items-center gap-1 text-xs">
						{#each locales as l (l)}
							{#if l === locale}
								<span class="font-semibold text-primary">{l.toUpperCase()}</span>
							{:else}
								<a
									href={localizeHref(page.url.pathname, { locale: l })}
									data-sveltekit-reload
									class="text-muted-foreground hover:text-primary"
								>
									{l.toUpperCase()}
								</a>
							{/if}
						{/each}
					</div>
					<span class="hidden text-sm text-muted-foreground sm:inline">
						{data.user?.fullName || data.user?.username}
					</span>
					<Button
						variant="ghost"
						size="icon-sm"
						aria-label={m.sidebar_logout()}
						disabled={loggingOut}
						onclick={handleLogout}
					>
						<LogOut class="size-4" />
					</Button>
				</div>
			</header>

			<Sheet.Root bind:open={drawerOpen}>
				<Sheet.Content side="left" class="w-64 gap-0 p-0">
					<Sheet.Header class="border-b border-border px-5 py-4 text-left">
						<Sheet.Title class="text-xl font-black text-primary">{m.sidebar_brand()}</Sheet.Title>
						<Sheet.Description class="text-xs text-muted-foreground">
							{m.sidebar_brandSub()}
						</Sheet.Description>
					</Sheet.Header>
					<nav class="flex-1 p-2">
						{#each navItems as item (item.href)}
							{@const active = isActive(item.href)}
							<a
								href={localizeHref(item.href)}
								onclick={() => (drawerOpen = false)}
								class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors {active
									? 'bg-muted font-semibold text-primary'
									: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
							>
								<item.icon class="size-4" />
								{item.label()}
							</a>
						{/each}
					</nav>
					<div class="border-t border-border p-3">
						<Button
							variant="ghost"
							size="sm"
							class="w-full justify-start"
							disabled={loggingOut}
							onclick={handleLogout}
						>
							<LogOut class="size-4" />
							{m.sidebar_logout()}
						</Button>
					</div>
				</Sheet.Content>
			</Sheet.Root>

			<main class="flex-1 p-4 lg:p-6">
				{@render children()}
			</main>
		</div>
	</div>
{/if}

<Toaster richColors position="top-center" />
