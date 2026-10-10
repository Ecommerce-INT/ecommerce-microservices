import { query } from '$app/server';
import { api } from '$lib/server/api';

/**
 * Dashboard aggregates. Counts come from the paginated endpoints
 * (`totalElements`), the lists from the newest-first pages. A real revenue
 * total needs an aggregate endpoint — see docs/web/backlog.md.
 */
export const getDashboard = query(async () => {
	const client = api();
	const [products, categories, orders] = await Promise.all([
		client.products.list({ page: 0, size: 5 }),
		client.categories.list({ page: 0, size: 1 }),
		client.orders.list({ page: 0, size: 5 })
	]);

	return {
		productCount: products.totalElements,
		categoryCount: categories.totalElements,
		orderCount: orders.totalElements,
		recentProducts: products.content,
		recentOrders: orders.content
	};
});
