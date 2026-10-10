import type {
	Category,
	CategoryPayload,
	CreateOrderPayload,
	CreatePaymentPayload,
	InventoryItem,
	Order,
	PaginatedResponse,
	Payment,
	Product,
	ProductListParams,
	ProductPayload,
	TokenResponse,
	UserResponse
} from '../types';
import { apiBaseUrl } from './env';

export class ApiError extends Error {
	readonly status: number;
	readonly body: unknown;

	constructor(status: number, message: string, body?: unknown) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.body = body;
	}
}

function query(params: Record<string, string | number | undefined>): string {
	const search = new URLSearchParams();
	for (const [key, value] of Object.entries(params)) {
		if (value !== undefined && value !== '') search.set(key, String(value));
	}
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

async function request<T>(
	base: string,
	path: string,
	init: RequestInit = {},
	token?: string | null
): Promise<T> {
	const headers = new Headers(init.headers);
	headers.set('accept', 'application/json');
	if (init.body && !headers.has('content-type')) headers.set('content-type', 'application/json');
	if (token) headers.set('authorization', `Bearer ${token}`);

	const response = await fetch(`${base}${path}`, { ...init, headers });
	const text = await response.text();

	let body: unknown;
	if (text) {
		try {
			body = JSON.parse(text);
		} catch {
			body = text;
		}
	}

	if (!response.ok) {
		let message = `${response.status} ${response.statusText}`;
		if (body && typeof body === 'object' && 'message' in body) {
			const candidate = (body as { message?: unknown }).message;
			if (typeof candidate === 'string' && candidate) message = candidate;
		}
		throw new ApiError(response.status, message, body);
	}

	// Unwrap the `ApiResponse<T>` envelope when present; raw payloads pass through.
	if (body && typeof body === 'object' && 'data' in body) {
		return (body as { data?: T }).data as T;
	}
	return body as T;
}

export interface Api {
	products: {
		list(params?: ProductListParams): Promise<PaginatedResponse<Product>>;
		get(id: number): Promise<Product>;
		create(payload: ProductPayload): Promise<Product>;
		update(id: number, payload: ProductPayload): Promise<Product>;
		remove(id: number): Promise<void>;
	};
	categories: {
		list(params?: { page?: number; size?: number }): Promise<PaginatedResponse<Category>>;
		get(id: number): Promise<Category>;
		create(payload: CategoryPayload): Promise<Category>;
		update(id: number, payload: CategoryPayload): Promise<Category>;
		remove(id: number): Promise<void>;
	};
	orders: {
		list(params?: { page?: number; size?: number }): Promise<PaginatedResponse<Order>>;
		get(id: number): Promise<Order>;
		create(payload: CreateOrderPayload): Promise<Order>;
	};
	payments: {
		create(payload: CreatePaymentPayload): Promise<Payment>;
	};
	inventory: {
		checkStock(skuCode: string): Promise<InventoryItem>;
	};
	auth: {
		exchangeTicket(ticket: string): Promise<TokenResponse>;
		refresh(refreshToken: string): Promise<TokenResponse>;
		logout(refreshToken: string): Promise<void>;
		profile(): Promise<UserResponse>;
	};
}

/**
 * Creates an API client bound to the gateway. Call per request with the
 * session's access token (see each app's `$lib/server/session.ts`).
 */
export function createApi(options: { token?: string | null; baseUrl?: string } = {}): Api {
	const base = options.baseUrl ?? apiBaseUrl();
	const token = options.token;
	const json = <T>(path: string, method: string, payload?: unknown) =>
		request<T>(
			base,
			path,
			{ method, body: payload === undefined ? undefined : JSON.stringify(payload) },
			token
		);

	return {
		products: {
			list: (params = {}) => request(base, `/api/products${query({ ...params })}`, {}, token),
			get: (id) => request(base, `/api/products/${id}`, {}, token),
			create: (payload) => json('/api/products', 'POST', payload),
			update: (id, payload) => json(`/api/products/${id}`, 'PUT', payload),
			remove: (id) => json(`/api/products/${id}`, 'DELETE')
		},
		categories: {
			list: (params = {}) => request(base, `/api/categories${query({ ...params })}`, {}, token),
			get: (id) => request(base, `/api/categories/${id}`, {}, token),
			create: (payload) => json('/api/categories', 'POST', payload),
			update: (id, payload) => json(`/api/categories/${id}`, 'PUT', payload),
			remove: (id) => json(`/api/categories/${id}`, 'DELETE')
		},
		orders: {
			list: (params = {}) => request(base, `/api/orders${query({ ...params })}`, {}, token),
			get: (id) => request(base, `/api/orders/${id}`, {}, token),
			create: (payload) => json('/api/orders', 'POST', payload)
		},
		payments: {
			create: (payload) => json('/api/payments', 'POST', payload)
		},
		inventory: {
			checkStock: (skuCode) => request(base, `/api/inventory${query({ skuCode })}`, {}, token)
		},
		auth: {
			exchangeTicket: (ticket) =>
				request(base, `/api/v1/auth/session${query({ ticket })}`, {}, null),
			refresh: (refreshToken) => json('/api/v1/auth/refresh', 'POST', { refreshToken }),
			logout: (refreshToken) => json('/api/v1/auth/logout', 'POST', { refreshToken }),
			profile: () => request(base, '/api/v1/users/me', {}, token)
		}
	};
}
