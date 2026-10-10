import { command, form, query, requested } from '$app/server';
import * as v from 'zod';
import { categoryUpsertSchema, pageQuerySchema } from '@ecommerce/lib/schemas';
import { api } from '$lib/server/api';

export const getCategories = query(pageQuerySchema, ({ page, size }) =>
	api().categories.list({ page, size })
);

export const saveCategory = form(categoryUpsertSchema, async ({ categoryId, ...input }) => {
	const client = api();
	const payload = { categoryTitle: input.categoryTitle, imageUrl: input.imageUrl || null };

	if (categoryId) {
		await client.categories.update(categoryId, payload);
	} else {
		await client.categories.create(payload);
	}

	return { saved: true };
});

export const deleteCategory = command(v.number().int().positive(), async (categoryId) => {
	await api().categories.remove(categoryId);
	await requested(getCategories, 5).refreshAll();
});
