import { redirect } from '@sveltejs/kit';
import { localizeHref } from '#lib/paraglide/runtime.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = () => {
	redirect(303, localizeHref('/dashboard'));
};
