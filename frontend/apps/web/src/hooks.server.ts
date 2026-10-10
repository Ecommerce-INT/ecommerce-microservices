import type { Handle } from '@sveltejs/kit/hooks';
import { sequence } from '@sveltejs/kit/hooks';
import { createApi } from '@ecommerce/lib/server';
import { getTextDirection } from '#lib/paraglide/runtime.js';
import { paraglideMiddleware } from '#lib/paraglide/server.js';
import {
	clearSession,
	refreshToken,
	session,
	setSession,
	shouldRefresh
} from '@ecommerce/lib/server';

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

/**
 * BFF session: exposes the identity from the httpOnly cookies to the app and
 * transparently rotates tokens when the access token expires. A failed refresh
 * clears the session (it never bounces back into SSO).
 */
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

export const handle = sequence(handleParaglide, handleSession);
