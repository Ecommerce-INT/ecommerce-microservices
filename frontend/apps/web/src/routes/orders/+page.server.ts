import { redirect } from '@sveltejs/kit';
import { deLocalizeUrl, localizeHref } from '$lib/paraglide/runtime.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals, url }) => {
	if (!locals.user) {
		const target = deLocalizeUrl(url).pathname + url.search;
		redirect(303, `${localizeHref('/login')}?redirect=${encodeURIComponent(target)}`);
	}
	return { user: locals.user };
};
