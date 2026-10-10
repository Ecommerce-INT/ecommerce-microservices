import { getRequestEvent } from '$app/server';
import { createApi, type Api } from '@ecommerce/lib/server';
import { accessToken } from '@ecommerce/lib/server';

/**
 * Request-scoped gateway client. Reads the session's access token from the
 * httpOnly cookie and points at `API_INTERNAL_URL`. Only usable inside a
 * request context (remote functions, load functions, hooks).
 */
export function api(): Api {
	const { cookies } = getRequestEvent();
	return createApi({ token: accessToken(cookies) });
}
