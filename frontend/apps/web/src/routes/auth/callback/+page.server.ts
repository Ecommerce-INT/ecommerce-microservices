import { redirect } from '@sveltejs/kit';
import { createApi } from '@ecommerce/lib/server';
import { localizeHref } from '$lib/paraglide/runtime.js';
import { safeRedirect } from '@ecommerce/lib/server';
import { clearSession, setSession, toSessionUser } from '@ecommerce/lib/server';
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

		// The profile supplies the roles the UI gates on. A failing profile
		// service must not block login: keep the session with no roles, which
		// fails closed for the backoffice guard.
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
