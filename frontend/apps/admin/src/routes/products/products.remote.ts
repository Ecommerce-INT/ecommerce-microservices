import { command, form, query, requested } from '$app/server';
import * as v from 'zod';
import { productQuerySchema, productUpsertSchema } from '@ecommerce/lib/schemas';
import { api } from '$lib/server/api';

export const getProducts = query(productQuerySchema, ({ page, size, search }) =>
	api().products.list({ page, size, search })
);

export const getCategoryOptions = query(async () => {
	const result = await api().categories.list({ page: 0, size: 100 });
	return result.content;
});

/** Create or update, depending on the optional `productId` field. */
export const saveProduct = form(productUpsertSchema, async ({ productId, ...input }) => {
	const client = api();
	const payload = {
		productTitle: input.productTitle,
		sku: input.sku,
		imageUrl: input.imageUrl || null,
		priceUnit: input.priceUnit,
		quantity: input.quantity,
		categoryId: input.categoryId ?? null
	};

	if (productId) {
		await client.products.update(productId, payload);
	} else {
		await client.products.create(payload);
	}

	return { saved: true };
});

export const deleteProduct = command(v.number().int().positive(), async (productId) => {
	await api().products.remove(productId);
	await requested(getProducts, 10).refreshAll();
});
