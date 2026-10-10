import { query } from '$app/server';
import { pageQuerySchema } from '@ecommerce/lib/schemas';
import { api } from '$lib/server/api';

export const getOrders = query(pageQuerySchema, ({ page, size }) =>
	api().orders.list({ page, size })
);
