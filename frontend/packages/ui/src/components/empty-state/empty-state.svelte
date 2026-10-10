<script lang="ts">
	import type { Component, Snippet } from 'svelte';
	import { cn } from '@ecommerce/ui/utils';

	let {
		icon,
		title,
		description,
		action,
		class: className
	}: {
		icon?: Component;
		title: string;
		description?: string;
		action?: Snippet;
		class?: string;
	} = $props();
</script>

<div
	class={cn(
		'border-border bg-card grid justify-items-center gap-4 rounded-xl border p-12 text-center',
		className
	)}
>
	{#if icon}
		{@const Icon = icon}
		<Icon class="text-muted-foreground size-12" />
	{/if}
	<div>
		<p class="text-lg font-semibold">{title}</p>
		{#if description}
			<p class="text-muted-foreground mt-1 text-sm">{description}</p>
		{/if}
	</div>
	{#if action}
		{@render action()}
	{/if}
</div>
