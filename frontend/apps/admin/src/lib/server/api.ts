import { getRequestEvent } from '$app/server';
import { accessToken, createApi, type Api } from '@ecommerce/lib/server';

/**
 * Request-scoped gateway client for the backoffice. Reads the access token from
 * the session cookie so every call carries the ADMIN user's JWT.
 */
export function api(): Api {
	const { cookies } = getRequestEvent();
	return createApi({ token: accessToken(cookies) });
}
