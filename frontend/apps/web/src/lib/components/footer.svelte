<script lang="ts">
	import { Mail, MapPin, Phone } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';
	import { localizeHref } from '#lib/paraglide/runtime.js';

	const columns = [
		{
			title: m.footer_categories,
			links: [
				{ href: '/products?categoryId=1', label: m.category_c1 },
				{ href: '/products?categoryId=2', label: m.category_c2 },
				{ href: '/products?categoryId=3', label: m.category_c3 },
				{ href: '/products?categoryId=4', label: m.category_c4 },
				{ href: '/products?categoryId=5', label: m.category_c5 }
			]
		},
		{
			title: m.footer_support,
			links: [
				{ href: '/products', label: m.header_allProducts },
				{ href: '/orders', label: m.footer_trackOrder },
				{ href: '/cart', label: m.header_cart },
				{ href: '/login', label: m.header_login }
			]
		}
	];
</script>

<footer class="mt-10 border-t border-border bg-card">
	<div class="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
		<div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
			<div>
				<p class="text-2xl font-black text-primary">EzBuy</p>
				<p class="mt-3 text-sm text-muted-foreground">{m.footer_tagline()}</p>
			</div>

			{#each columns as column (column.title)}
				<div>
					<h3 class="mb-3 text-sm font-semibold">{column.title()}</h3>
					<ul class="space-y-2 text-sm text-muted-foreground">
						{#each column.links as link (link.href + link.label())}
							<li>
								<a href={localizeHref(link.href)} class="hover:text-primary">{link.label()}</a>
							</li>
						{/each}
					</ul>
				</div>
			{/each}

			<div>
				<h3 class="mb-3 text-sm font-semibold">{m.footer_contactTitle()}</h3>
				<ul class="space-y-2 text-sm text-muted-foreground">
					<li class="flex items-start gap-2">
						<MapPin class="mt-0.5 size-4 shrink-0" />
						{m.footer_address()}
					</li>
					<li class="flex items-center gap-2">
						<Phone class="size-4 shrink-0" />
						{m.footer_phone()}
					</li>
					<li class="flex items-center gap-2">
						<Mail class="size-4 shrink-0" />
						{m.footer_email()}
					</li>
				</ul>
			</div>
		</div>

		<div
			class="mt-8 flex flex-col gap-2 border-t border-border pt-6 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between"
		>
			<span>{m.footer_copyright()}</span>
			<span>{m.footer_paymentTitle()}</span>
		</div>
	</div>
</footer>
