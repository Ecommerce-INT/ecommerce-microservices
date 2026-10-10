import type { Cookies } from '@sveltejs/kit';
import type { Role, SessionUser, TokenResponse, UserResponse } from '../types';

export const ACCESS_COOKIE = 'access_token';
export const REFRESH_COOKIE = 'refresh_token';
export const USER_COOKIE = 'session_user';
export const EXPIRES_COOKIE = 'session_expires';

const base = { path: '/', httpOnly: true, sameSite: 'lax' } as const;
const ACCESS_FALLBACK_SECONDS = 300;
const REFRESH_FALLBACK_SECONDS = 1800;

export function toSessionUser(profile: UserResponse): SessionUser {
	return {
		id: profile.id,
		username: profile.username,
		fullName: profile.fullName,
		email: profile.email,
		roles: (profile.roles ?? []).map((role: Role) => role.name)
	};
}

export function isAdmin(user: SessionUser | null | undefined): boolean {
	return !!user?.roles.some((role) => role === 'ADMIN' || role === 'ROLE_ADMIN');
}

/**
 * Persists a token bundle (from the SSO ticket exchange or a refresh) plus the
 * profile. `user: undefined` keeps the existing profile cookie — refreshes
 * rotate tokens but not identity.
 */
export function setSession(cookies: Cookies, tokens: TokenResponse, user?: SessionUser | null) {
	const accessMaxAge = tokens.expires_in > 0 ? tokens.expires_in : ACCESS_FALLBACK_SECONDS;
	const refreshMaxAge =
		tokens.refresh_expires_in > 0 ? tokens.refresh_expires_in : REFRESH_FALLBACK_SECONDS;

	cookies.set(ACCESS_COOKIE, tokens.access_token, { ...base, maxAge: accessMaxAge });
	cookies.set(REFRESH_COOKIE, tokens.refresh_token, { ...base, maxAge: refreshMaxAge });
	cookies.set(EXPIRES_COOKIE, String(Math.floor(Date.now() / 1000) + accessMaxAge), {
		...base,
		maxAge: refreshMaxAge
	});

	if (user === null) {
		cookies.delete(USER_COOKIE, { path: '/' });
	} else if (user) {
		cookies.set(USER_COOKIE, JSON.stringify(user), { ...base, maxAge: refreshMaxAge });
	}
}

export function clearSession(cookies: Cookies) {
	for (const name of [ACCESS_COOKIE, REFRESH_COOKIE, USER_COOKIE, EXPIRES_COOKIE]) {
		cookies.delete(name, { path: '/' });
	}
}

export function accessToken(cookies: Cookies): string | null {
	return cookies.get(ACCESS_COOKIE) ?? null;
}

export function refreshToken(cookies: Cookies): string | null {
	return cookies.get(REFRESH_COOKIE) ?? null;
}

export interface Session {
	user: SessionUser;
}

/** Reads the session cookie pair; null when there is no usable identity. */
export function session(cookies: Cookies): Session | null {
	const raw = cookies.get(USER_COOKIE);
	if (!raw) return null;
	try {
		const user = JSON.parse(raw) as SessionUser;
		if (!user || typeof user.username !== 'string') return null;
		return { user };
	} catch {
		return null;
	}
}

/**
 * True when the access token is gone or within the refresh window. The refresh
 * token cookie outlives the access token cookie by design, so a missing access
 * token is the normal "needs refresh" signal.
 */
export function shouldRefresh(cookies: Cookies): boolean {
	if (!accessToken(cookies)) return true;
	const expiresAt = Number(cookies.get(EXPIRES_COOKIE) ?? 0);
	return !expiresAt || expiresAt - Math.floor(Date.now() / 1000) < 30;
}
