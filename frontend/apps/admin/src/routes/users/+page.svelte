<script lang="ts">
	import { m } from '#lib/paraglide/messages.js';
	import ComingSoon from '#lib/components/coming-soon.svelte';
</script>

<svelte:head>
	<title>{m.meta_template({ page: m.users_metaTitle() })}</title>
</svelte:head>

<ComingSoon title={m.users_title()} />
