import { describe, expect, test } from 'bun:test';
import {
	checkoutSchema,
	pageQuerySchema,
	productFormSchema,
	productQuerySchema,
	productUpsertSchema,
	shippingSchema
} from './schemas';

const validShipping = {
	fullName: 'Nguyen Van An',
	phone: '0901234567',
	address: '12 Le Loi',
	city: 'Ho Chi Minh',
	district: 'District 1'
};

describe('pageQuerySchema', () => {
	test('applies defaults', () => {
		expect(pageQuerySchema.parse({})).toEqual({ page: 0, size: 12 });
	});

	test('coerces numeric strings', () => {
		expect(pageQuerySchema.parse({ page: '2', size: '24' })).toEqual({ page: 2, size: 24 });
	});

	test('rejects out-of-range values', () => {
		expect(pageQuerySchema.safeParse({ page: -1 }).success).toBe(false);
		expect(pageQuerySchema.safeParse({ size: 0 }).success).toBe(false);
		expect(pageQuerySchema.safeParse({ size: 61 }).success).toBe(false);
	});
});

describe('productQuerySchema', () => {
	test('accepts only supported sort values', () => {
		expect(productQuerySchema.safeParse({ sort: 'priceUnit,desc' }).success).toBe(true);
		expect(productQuerySchema.safeParse({ sort: 'productTitle,asc' }).success).toBe(false);
	});

	test('coerces categoryId and drops empty strings', () => {
		expect(productQuerySchema.parse({ categoryId: '42' }).categoryId).toBe(42);
		expect(productQuerySchema.parse({ categoryId: 7 }).categoryId).toBe(7);
		expect(productQuerySchema.parse({ categoryId: '' }).categoryId).toBeUndefined();
		expect(productQuerySchema.parse({}).categoryId).toBeUndefined();
	});

	test('rejects invalid categoryId and whitespace-only search', () => {
		expect(productQuerySchema.safeParse({ categoryId: 'abc' }).success).toBe(false);
		expect(productQuerySchema.safeParse({ categoryId: '-3' }).success).toBe(false);
		expect(productQuerySchema.safeParse({ search: '   ' }).success).toBe(false);
	});

	test('trims search terms', () => {
		expect(productQuerySchema.parse({ search: '  shoe  ' }).search).toBe('shoe');
	});
});

describe('productFormSchema', () => {
	const valid = {
		productTitle: 'Sneaker',
		sku: 'SNK-1',
		imageUrl: '',
		priceUnit: '250000',
		quantity: '3'
	};

	test('coerces numeric form fields', () => {
		const parsed = productFormSchema.parse(valid);
		expect(parsed.priceUnit).toBe(250000);
		expect(parsed.quantity).toBe(3);
		expect(parsed.categoryId).toBeUndefined();
	});

	test('rejects non-numeric and empty numeric fields with the custom message', () => {
		const result = productFormSchema.safeParse({ ...valid, priceUnit: 'abc' });
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.issues[0]?.message).toBe('Please enter a price');
		}
		expect(productFormSchema.safeParse({ ...valid, quantity: '' }).success).toBe(false);
	});

	test('rejects negative prices and fractional quantities', () => {
		expect(productFormSchema.safeParse({ ...valid, priceUnit: '-5' }).success).toBe(false);
		expect(productFormSchema.safeParse({ ...valid, quantity: '2.5' }).success).toBe(false);
	});
});

describe('productUpsertSchema', () => {
	const base = {
		productTitle: 'Sneaker',
		sku: 'SNK-1',
		imageUrl: '',
		priceUnit: '1',
		quantity: '1'
	};

	test('treats a missing or empty id as create', () => {
		expect(productUpsertSchema.parse(base).productId).toBeUndefined();
		expect(productUpsertSchema.parse({ ...base, productId: '' }).productId).toBeUndefined();
	});

	test('coerces a provided id to a number', () => {
		expect(productUpsertSchema.parse({ ...base, productId: '7' }).productId).toBe(7);
	});
});

describe('shippingSchema', () => {
	test('accepts a valid payload and trims strings', () => {
		const parsed = shippingSchema.parse({ ...validShipping, fullName: '  Nguyen Van An  ' });
		expect(parsed.fullName).toBe('Nguyen Van An');
	});

	test('rejects malformed phones and short names/addresses', () => {
		expect(shippingSchema.safeParse({ ...validShipping, phone: 'abc' }).success).toBe(false);
		expect(shippingSchema.safeParse({ ...validShipping, fullName: ' A ' }).success).toBe(false);
		expect(shippingSchema.safeParse({ ...validShipping, address: 'abc' }).success).toBe(false);
		expect(shippingSchema.safeParse({ ...validShipping, city: '' }).success).toBe(false);
	});

	test('caps the note length', () => {
		expect(shippingSchema.safeParse({ ...validShipping, note: 'x'.repeat(501) }).success).toBe(
			false
		);
		expect(shippingSchema.safeParse({ ...validShipping, note: 'ok' }).success).toBe(true);
	});
});

describe('checkoutSchema', () => {
	const base = { shipping: validShipping, paymentMethod: 'COD' as const };

	test('requires at least one line item', () => {
		expect(checkoutSchema.safeParse({ ...base, items: [] }).success).toBe(false);
	});

	test('bounds item quantity between 1 and 99', () => {
		expect(
			checkoutSchema.safeParse({ ...base, items: [{ productId: 1, quantity: 0 }] }).success
		).toBe(false);
		expect(
			checkoutSchema.safeParse({ ...base, items: [{ productId: 1, quantity: 100 }] }).success
		).toBe(false);
		expect(
			checkoutSchema.safeParse({ ...base, items: [{ productId: 1, quantity: 99 }] }).success
		).toBe(true);
	});

	test('accepts only COD and ONLINE payment methods', () => {
		expect(checkoutSchema.safeParse({ ...base, paymentMethod: 'BITCOIN', items: [] }).success).toBe(
			false
		);
		expect(checkoutSchema.safeParse({ ...base, paymentMethod: 'ONLINE', items: [] }).success).toBe(
			false
		);
	});
});
