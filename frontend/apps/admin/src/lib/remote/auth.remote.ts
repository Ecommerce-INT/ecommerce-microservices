import { command, getRequestEvent } from '$app/server';
import { createApi } from '@ecommerce/lib/server';
import { clearSession, refreshToken } from '@ecommerce/lib/server';

/** Ends the gateway session (best effort) and clears the BFF cookies. */
export const logout = command(async () => {
	const { cookies } = getRequestEvent();
	const token = refreshToken(cookies);

	if (token) {
		try {
			await createApi().auth.logout(token);
		} catch {
			// The local session is cleared regardless of the gateway's answer.
		}
	}

	clearSession(cookies);
});
