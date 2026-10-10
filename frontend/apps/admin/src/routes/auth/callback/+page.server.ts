import { redirect } from '@sveltejs/kit';
import { createApi, safeRedirect, setSession, toSessionUser } from '@ecommerce/lib/server';
import { clearSession } from '@ecommerce/lib/server';
import { localizeHref } from '#lib/paraglide/runtime.js';
import type { PageServerLoad } from './$types';

const REDIRECT_COOKIE = 'auth_redirect';

export const load: PageServerLoad = async ({ url, cookies }) => {
	const target = safeRedirect(cookies.get(REDIRECT_COOKIE));
	cookies.delete(REDIRECT_COOKIE, { path: '/' });

	const ticket = url.searchParams.get('ticket');
	if (!ticket || url.searchParams.has('error')) {
		clearSession(cookies);
		return { failed: true };
	}

	try {
		const tokens = await createApi().auth.exchangeTicket(ticket);
		let user = null;
		try {
			user = toSessionUser(await createApi({ token: tokens.access_token }).auth.profile());
		} catch {
			user = null;
		}
		setSession(cookies, tokens, user);
	} catch {
		clearSession(cookies);
		return { failed: true };
	}

	redirect(303, localizeHref(target));
};
