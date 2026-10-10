import { query } from '$app/server';
import { productQuerySchema } from '@ecommerce/lib/schemas';
import { api } from '$lib/server/api';

/**
 * Catalog listing. `productQuerySchema` coerces/validates every argument that
 * arrives from the client (numbers from URL params, optional filters).
 *
 * NOTE: the `search` filter is passed through to the gateway; whether
 * product-service honours it server-side is verified against the running
 * stack (the old app never sent it at all).
 */
export const getProducts = query(productQuerySchema, (input) => api().products.list(input));

/** Categories for the filter sidebar (first 50 — good enough for the catalog). */
export const getCategories = query(async () => {
	const result = await api().categories.list({ page: 0, size: 50 });
	return result.content;
});
