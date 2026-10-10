import { redirect } from '@sveltejs/kit';
import { apiPublicUrl } from '@ecommerce/lib/server';
import { localizeHref } from '#lib/paraglide/runtime.js';
import { safeRedirect } from '@ecommerce/lib/server';
import type { PageServerLoad } from './$types';

const REDIRECT_COOKIE = 'auth_redirect';

export const load: PageServerLoad = ({ url, cookies, locals }) => {
	const target = safeRedirect(url.searchParams.get('redirect'));

	if (locals.user) redirect(303, localizeHref(safeRedirect(url.searchParams.get('redirect'))));

	// The gateway only allows a fixed `/auth/callback` redirect, so the final
	// destination rides in a short-lived cookie instead of a query param.
	cookies.set(REDIRECT_COOKIE, target, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		maxAge: 600
	});

	const redirectUri = `${url.origin}/auth/callback`;
	const ssoUrl = `${apiPublicUrl()}/api/v1/auth/login?redirect_uri=${encodeURIComponent(redirectUri)}`;

	return { ssoUrl };
};
