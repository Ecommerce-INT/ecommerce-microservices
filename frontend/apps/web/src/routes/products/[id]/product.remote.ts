import { query } from '$app/server';
import * as v from 'zod';
import { api } from '$lib/server/api';

export const getProduct = query(v.coerce.number().int().positive(), (id) => api().products.get(id));
