import { query } from '$app/server';
import { pageQuerySchema } from '@ecommerce/lib/schemas';
import { api } from '#lib/server/api';

/** Orders for the signed-in user (the gateway scopes them by the JWT). */
export const getOrders = query(pageQuerySchema, ({ page, size }) =>
	api().orders.list({ page, size })
);
