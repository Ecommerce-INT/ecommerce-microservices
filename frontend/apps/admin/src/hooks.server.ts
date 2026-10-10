import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit/hooks';
import { sequence } from '@sveltejs/kit/hooks';
import {
	clearSession,
	createApi,
	isAdmin,
	refreshToken,
	session,
	setSession,
	shouldRefresh
} from '@ecommerce/lib/server';
import { deLocalizeUrl, getTextDirection, localizeHref } from '#lib/paraglide/runtime.js';
import { paraglideMiddleware } from '#lib/paraglide/server.js';

/** Routes reachable without an admin session (everything else is guarded). */
const PUBLIC_ROUTES = ['/login', '/auth/callback', '/access-denied'];

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		return resolve(
			{ ...event, request },
			{
				transformPageChunk: ({ html }) =>
					html
						.replace('%paraglide.lang%', locale)
						.replace('%paraglide.dir%', getTextDirection(locale))
			}
		);
	});

const handleSession: Handle = async ({ event, resolve }) => {
	let current = session(event.cookies);

	if (current && shouldRefresh(event.cookies)) {
		const token = refreshToken(event.cookies);
		if (token) {
			try {
				const tokens = await createApi().auth.refresh(token);
				setSession(event.cookies, tokens);
				current = session(event.cookies);
			} catch {
				clearSession(event.cookies);
				current = null;
			}
		}
	}

	event.locals.user = current?.user ?? null;
	return resolve(event);
};

/**
 * Strict backoffice guard: every route except the SSO entry points requires a
 * signed-in user whose token carries an ADMIN role. Anonymous visitors are sent
 * to /login with a return path; signed-in non-admins get a 403.
 */
const handleAdminGuard: Handle = ({ event, resolve }) => {
	const { pathname } = deLocalizeUrl(event.url);
	const isPublic =
		PUBLIC_ROUTES.some((route) => pathname === route || pathname.startsWith(`${route}/`)) ||
		pathname.startsWith('/_app/') ||
		pathname.includes('.');

	if (isPublic) return resolve(event);

	if (!event.locals.user) {
		const target = pathname + event.url.search;
		redirect(303, `${localizeHref('/login')}?redirect=${encodeURIComponent(target)}`);
	}

	if (!isAdmin(event.locals.user)) {
		redirect(303, localizeHref('/access-denied'));
	}

	return resolve(event);
};

export const handle = sequence(handleParaglide, handleSession, handleAdminGuard);
