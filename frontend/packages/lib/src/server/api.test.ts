import { afterEach, describe, expect, test } from 'bun:test';
import { ApiError, createApi } from './api';

const originalFetch = globalThis.fetch;

interface FetchCall {
	url: string;
	init?: RequestInit;
}

function stubFetch(handler: (call: FetchCall) => Response | Promise<Response>) {
	const calls: FetchCall[] = [];
	globalThis.fetch = (async (input: string | URL | Request, init?: RequestInit) => {
		const call = { url: String(input), init };
		calls.push(call);
		return handler(call);
	}) as typeof fetch;
	return calls;
}

function json(body: unknown, init: ResponseInit = {}) {
	return new Response(JSON.stringify(body), {
		status: 200,
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function headersOf(call: FetchCall) {
	return new Headers(call.init?.headers);
}

afterEach(() => {
	globalThis.fetch = originalFetch;
});

describe('createApi', () => {
	test('unwraps the ApiResponse envelope', async () => {
		stubFetch(() =>
			json({
				success: true,
				code: 'OK',
				timestamp: '2026-01-01T00:00:00Z',
				data: { productId: 1, productTitle: 'Sneaker', priceUnit: 100, quantity: 2 }
			})
		);
		const api = createApi({ baseUrl: 'http://gateway' });
		const product = await api.products.get(1);
		expect(product.productId).toBe(1);
		expect(product.productTitle).toBe('Sneaker');
	});

	test('passes raw payloads through untouched', async () => {
		stubFetch(() => json({ productId: 2, productTitle: 'Boot', priceUnit: 50, quantity: 1 }));
		const api = createApi({ baseUrl: 'http://gateway' });
		await expect(api.products.get(2)).resolves.toMatchObject({ productId: 2 });
	});

	test('returns undefined for empty bodies', async () => {
		stubFetch(() => new Response(null, { status: 204 }));
		const api = createApi({ baseUrl: 'http://gateway' });
		await expect(api.products.remove(3)).resolves.toBeUndefined();
	});

	test('builds query strings and skips empty parameters', async () => {
		const calls = stubFetch(() => json({ content: [], totalElements: 0 }));
		const api = createApi({ baseUrl: 'http://gateway' });
		await api.products.list({ page: 1, size: 12, search: '', categoryId: undefined });
		expect(calls[0]?.url).toBe('http://gateway/api/products?page=1&size=12');
	});

	test('attaches the bearer token when provided', async () => {
		const calls = stubFetch(() => json({ data: {} }));
		const api = createApi({ baseUrl: 'http://gateway', token: 'secret-token' });
		await api.products.get(1);
		expect(calls).toHaveLength(1);
		expect(headersOf(calls[0]!).get('authorization')).toBe('Bearer secret-token');
	});

	test('omits the authorization header without a token', async () => {
		const calls = stubFetch(() => json({ data: {} }));
		const api = createApi({ baseUrl: 'http://gateway' });
		await api.products.get(1);
		expect(calls).toHaveLength(1);
		expect(headersOf(calls[0]!).get('authorization')).toBeNull();
	});

	test('sends JSON bodies on create', async () => {
		const calls = stubFetch(() => json({ data: { productId: 9 } }));
		const api = createApi({ baseUrl: 'http://gateway', token: 't' });
		await api.products.create({
			productTitle: 'Hat',
			sku: 'HAT-1',
			priceUnit: 100,
			quantity: 1
		});
		expect(calls).toHaveLength(1);
		const call = calls[0]!;
		expect(call.url).toBe('http://gateway/api/products');
		expect(call.init?.method).toBe('POST');
		expect(headersOf(call).get('content-type')).toBe('application/json');
		expect(JSON.parse(String(call.init?.body))).toEqual({
			productTitle: 'Hat',
			sku: 'HAT-1',
			priceUnit: 100,
			quantity: 1
		});
	});

	test('throws ApiError carrying the status and server message', async () => {
		stubFetch(() =>
			json({ message: 'Product not found' }, { status: 404, statusText: 'Not Found' })
		);
		const api = createApi({ baseUrl: 'http://gateway' });
		const error = await api.products.get(404).catch((caught: unknown) => caught);
		expect(error).toBeInstanceOf(ApiError);
		expect((error as ApiError).status).toBe(404);
		expect((error as ApiError).message).toBe('Product not found');
	});

	test('falls back to the HTTP status when the error body has no message', async () => {
		stubFetch(() => new Response('boom', { status: 500, statusText: 'Server Error' }));
		const api = createApi({ baseUrl: 'http://gateway' });
		const error = await api.products.get(1).catch((caught: unknown) => caught);
		expect((error as ApiError).message).toBe('500 Server Error');
	});

	test('exchanges the SSO ticket without a bearer token', async () => {
		const calls = stubFetch(() => json({ data: { access_token: 'a', refresh_token: 'r' } }));
		const api = createApi({ baseUrl: 'http://gateway', token: 'ignored' });
		await api.auth.exchangeTicket('ticket-1');
		expect(calls).toHaveLength(1);
		expect(calls[0]?.url).toBe('http://gateway/api/v1/auth/session?ticket=ticket-1');
		expect(headersOf(calls[0]!).get('authorization')).toBeNull();
	});
});
