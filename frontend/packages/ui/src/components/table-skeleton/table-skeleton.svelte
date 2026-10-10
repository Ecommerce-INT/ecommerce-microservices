<script lang="ts">
	import { Skeleton } from '@ecommerce/ui/skeleton';
	import { cn } from '@ecommerce/ui/utils';

	let { rows = 5, class: className }: { rows?: number; class?: string } = $props();
</script>

<div class={cn('border-border bg-card overflow-hidden rounded-xl border', className)}>
	{#each Array(rows).keys() as row (row)}
		<div class="border-border flex items-center gap-4 border-b p-4 last:border-b-0">
			<Skeleton class="size-9 shrink-0 rounded-lg" />
			<Skeleton class="h-4 flex-1" />
			<Skeleton class="hidden h-4 w-24 sm:block" />
			<Skeleton class="h-4 w-16" />
		</div>
	{/each}
</div>
