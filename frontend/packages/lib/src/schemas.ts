import * as v from 'zod';

/**
 * A numeric form field: arrives as a string from the DOM, leaves as a number.
 * Coerced schemas (`v.coerce.number()`) are not usable in remote forms because
 * their input type (`unknown`) fails SvelteKit's `RemoteFormInput` constraint.
 */
const formNumber = (message: string) =>
	v
		.string()
		.trim()
		.min(1, message)
		.regex(/^\d+$/, message)
		.transform((value) => Number(value));

/**
 * Coerces a single optional numeric value (URL params, select values, hidden
 * form fields) to a number. Deliberately narrow — remote forms only allow
 * string | number | boolean | File | undefined in their inputs.
 */
const optionalId = v
	.union([v.literal(''), v.string().regex(/^\d+$/), v.number().int().positive()])
	.optional()
	.transform((value) =>
		typeof value === 'string' && value !== ''
			? Number(value)
			: typeof value === 'number'
				? value
				: undefined
	);

export const pageQuerySchema = v.object({
	page: v.coerce.number().int().min(0).default(0),
	size: v.coerce.number().int().min(1).max(60).default(12)
});

export const productQuerySchema = pageQuerySchema.extend({
	sort: v.enum(['productId,asc', 'productId,desc', 'priceUnit,asc', 'priceUnit,desc']).optional(),
	categoryId: optionalId,
	search: v.string().trim().min(1).optional()
});

export const shippingSchema = v.object({
	fullName: v.string().trim().min(2, 'Please enter your full name'),
	phone: v
		.string()
		.trim()
		.regex(/^[0-9+\-\s]{8,15}$/, 'Please enter a valid phone number'),
	address: v.string().trim().min(5, 'Please enter a valid address'),
	city: v.string().trim().min(1, 'Please select a city'),
	district: v.string().trim().min(1, 'Please enter a district'),
	note: v.string().trim().max(500).optional()
});

export const checkoutLineSchema = v.object({
	productId: v.number().int().positive(),
	quantity: v.number().int().min(1).max(99)
});

export const checkoutSchema = v.object({
	shipping: shippingSchema,
	paymentMethod: v.enum(['COD', 'ONLINE']),
	items: v.array(checkoutLineSchema).min(1, 'Your cart is empty')
});

export const productFormSchema = v.object({
	productTitle: v.string().trim().min(1, 'Please enter a product name').max(200),
	sku: v.string().trim().min(1, 'Please enter an SKU').max(64),
	imageUrl: v.string().trim(),
	categoryId: optionalId,
	priceUnit: formNumber('Please enter a price'),
	quantity: formNumber('Please enter a quantity')
});

export const categoryFormSchema = v.object({
	categoryTitle: v.string().trim().min(1, 'Please enter a category name').max(120),
	imageUrl: v.string().trim()
});

/** Backoffice upsert schemas: same payloads plus the optional entity id. */
export const productUpsertSchema = productFormSchema.extend({ productId: optionalId });
export const categoryUpsertSchema = categoryFormSchema.extend({ categoryId: optionalId });

export type PageQuery = v.infer<typeof pageQuerySchema>;
export type ProductQuery = v.infer<typeof productQuerySchema>;
export type ShippingInput = v.infer<typeof shippingSchema>;
export type CheckoutInput = v.infer<typeof checkoutSchema>;
export type ProductFormInput = v.infer<typeof productFormSchema>;
export type CategoryFormInput = v.infer<typeof categoryFormSchema>;
export type ProductUpsertInput = v.infer<typeof productUpsertSchema>;
export type CategoryUpsertInput = v.infer<typeof categoryUpsertSchema>;
