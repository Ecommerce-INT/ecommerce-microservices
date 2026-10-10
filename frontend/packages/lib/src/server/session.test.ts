import type { Cookies } from '@sveltejs/kit';
import { describe, expect, test } from 'bun:test';
import type { SessionUser, TokenResponse, UserResponse } from '../types';
import {
	ACCESS_COOKIE,
	EXPIRES_COOKIE,
	REFRESH_COOKIE,
	USER_COOKIE,
	accessToken,
	clearSession,
	isAdmin,
	refreshToken,
	session,
	setSession,
	shouldRefresh,
	toSessionUser
} from './session';

interface SetRecord {
	name: string;
	value: string;
	options?: Record<string, unknown>;
}

function fakeCookies(initial: Record<string, string> = {}) {
	const store = new Map(Object.entries(initial));
	const sets: SetRecord[] = [];
	const deletes: string[] = [];
	const cookies = {
		get: (name: string) => store.get(name),
		set: (name: string, value: string, options?: Record<string, unknown>) => {
			sets.push({ name, value, options });
			store.set(name, value);
		},
		delete: (name: string) => {
			deletes.push(name);
			store.delete(name);
		}
	};
	return { cookies: cookies as unknown as Cookies, sets, deletes, store };
}

const tokens: TokenResponse = {
	access_token: 'access-1',
	refresh_token: 'refresh-1',
	expires_in: 300,
	refresh_expires_in: 1800,
	token_type: 'Bearer',
	scope: 'openid'
};

const user: SessionUser = {
	id: 1,
	username: 'an',
	fullName: 'Nguyen Van An',
	email: 'an@example.com',
	roles: ['USER']
};

describe('toSessionUser', () => {
	test('maps role objects to role names', () => {
		const profile: UserResponse = {
			id: 7,
			username: 'admin',
			fullName: 'Root',
			email: 'root@example.com',
			gender: 'M',
			roles: [
				{ id: 1, name: 'ADMIN' },
				{ id: 2, name: 'USER' }
			]
		};
		expect(toSessionUser(profile)).toEqual({
			id: 7,
			username: 'admin',
			fullName: 'Root',
			email: 'root@example.com',
			roles: ['ADMIN', 'USER']
		});
	});
});

describe('isAdmin', () => {
	test('accepts both role spellings', () => {
		expect(isAdmin({ ...user, roles: ['ADMIN'] })).toBe(true);
		expect(isAdmin({ ...user, roles: ['ROLE_ADMIN'] })).toBe(true);
	});

	test('rejects regular users and missing sessions', () => {
		expect(isAdmin({ ...user, roles: ['USER'] })).toBe(false);
		expect(isAdmin(null)).toBe(false);
		expect(isAdmin(undefined)).toBe(false);
	});
});

describe('setSession', () => {
	test('persists tokens, expiry and profile', () => {
		const { cookies, sets } = fakeCookies();
		setSession(cookies, tokens, user);

		const access = sets.find((entry) => entry.name === ACCESS_COOKIE);
		expect(access?.value).toBe('access-1');
		expect(access?.options).toMatchObject({
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 300
		});

		const refresh = sets.find((entry) => entry.name === REFRESH_COOKIE);
		expect(refresh?.value).toBe('refresh-1');
		expect(refresh?.options).toMatchObject({ maxAge: 1800 });

		const expires = sets.find((entry) => entry.name === EXPIRES_COOKIE);
		const expected = Math.floor(Date.now() / 1000) + 300;
		expect(Math.abs(Number(expires?.value) - expected)).toBeLessThanOrEqual(2);

		const profile = sets.find((entry) => entry.name === USER_COOKIE);
		expect(JSON.parse(profile?.value ?? '')).toEqual(user);
	});

	test('keeps the existing profile when user is omitted', () => {
		const { cookies, sets, deletes } = fakeCookies({ [USER_COOKIE]: JSON.stringify(user) });
		setSession(cookies, tokens);
		expect(sets.map((entry) => entry.name)).toEqual([
			ACCESS_COOKIE,
			REFRESH_COOKIE,
			EXPIRES_COOKIE
		]);
		expect(deletes).toHaveLength(0);
	});

	test('deletes the profile when user is null', () => {
		const { cookies, deletes } = fakeCookies({ [USER_COOKIE]: JSON.stringify(user) });
		setSession(cookies, tokens, null);
		expect(deletes).toEqual([USER_COOKIE]);
	});

	test('falls back to default max ages for non-positive expiries', () => {
		const { cookies, sets } = fakeCookies();
		setSession(cookies, { ...tokens, expires_in: 0, refresh_expires_in: 0 }, user);
		expect(sets.find((entry) => entry.name === ACCESS_COOKIE)?.options).toMatchObject({
			maxAge: 300
		});
		expect(sets.find((entry) => entry.name === REFRESH_COOKIE)?.options).toMatchObject({
			maxAge: 1800
		});
	});
});

describe('clearSession', () => {
	test('deletes every session cookie', () => {
		const { cookies, deletes } = fakeCookies({
			[ACCESS_COOKIE]: 'a',
			[REFRESH_COOKIE]: 'r',
			[USER_COOKIE]: JSON.stringify(user),
			[EXPIRES_COOKIE]: '123'
		});
		clearSession(cookies);
		const expected = [ACCESS_COOKIE, EXPIRES_COOKIE, REFRESH_COOKIE, USER_COOKIE];
		expect([...deletes].sort()).toEqual([...expected].sort());
	});
});

describe('session', () => {
	test('returns null without a profile cookie', () => {
		const { cookies } = fakeCookies();
		expect(session(cookies)).toBeNull();
	});

	test('returns null for malformed or incomplete payloads', () => {
		expect(session(fakeCookies({ [USER_COOKIE]: '{nope' }).cookies)).toBeNull();
		expect(session(fakeCookies({ [USER_COOKIE]: 'null' }).cookies)).toBeNull();
		expect(session(fakeCookies({ [USER_COOKIE]: '{"id":1}' }).cookies)).toBeNull();
	});

	test('deserializes a valid profile', () => {
		const { cookies } = fakeCookies({ [USER_COOKIE]: JSON.stringify(user) });
		expect(session(cookies)).toEqual({ user });
	});
});

describe('shouldRefresh', () => {
	test('is true without an access token', () => {
		const { cookies } = fakeCookies();
		expect(shouldRefresh(cookies)).toBe(true);
	});

	test('is true when the access token expires within the refresh margin', () => {
		const soon = Math.floor(Date.now() / 1000) + 10;
		const { cookies } = fakeCookies({ [ACCESS_COOKIE]: 'a', [EXPIRES_COOKIE]: String(soon) });
		expect(shouldRefresh(cookies)).toBe(true);
	});

	test('is true when the expiry cookie is missing', () => {
		const { cookies } = fakeCookies({ [ACCESS_COOKIE]: 'a' });
		expect(shouldRefresh(cookies)).toBe(true);
	});

	test('is false while the access token is still fresh', () => {
		const later = Math.floor(Date.now() / 1000) + 3600;
		const { cookies } = fakeCookies({ [ACCESS_COOKIE]: 'a', [EXPIRES_COOKIE]: String(later) });
		expect(shouldRefresh(cookies)).toBe(false);
	});
});

describe('token accessors', () => {
	test('read cookie values', () => {
		const { cookies } = fakeCookies({ [ACCESS_COOKIE]: 'a', [REFRESH_COOKIE]: 'r' });
		expect(accessToken(cookies)).toBe('a');
		expect(refreshToken(cookies)).toBe('r');
	});

	test('return null when unset', () => {
		const { cookies } = fakeCookies();
		expect(accessToken(cookies)).toBeNull();
		expect(refreshToken(cookies)).toBeNull();
	});
});
