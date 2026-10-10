<script lang="ts">
	import { Construction } from '@lucide/svelte';
	import { m } from '#lib/paraglide/messages.js';

	let { title }: { title: string } = $props();
</script>

<div
	class="grid justify-items-center gap-3 rounded-xl border border-border bg-card p-16 text-center"
>
	<Construction class="size-12 text-muted-foreground" />
	<h1 class="text-lg font-semibold">{title}</h1>
	<p class="text-sm text-muted-foreground">{m.common_developing()}</p>
</div>
