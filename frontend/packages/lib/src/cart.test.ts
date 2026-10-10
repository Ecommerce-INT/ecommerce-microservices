import { describe, expect, test } from 'bun:test';
import { computeTotals, FREE_SHIPPING_THRESHOLD, SHIPPING_FEE } from './cart';

describe('computeTotals', () => {
	test('an empty cart has no subtotal and no fee', () => {
		expect(computeTotals([])).toEqual({ subtotal: 0, shippingFee: 0, total: 0 });
	});

	test('sums lines and charges a flat fee below the threshold', () => {
		expect(computeTotals([{ priceUnit: 100_000, quantity: 2 }])).toEqual({
			subtotal: 200_000,
			shippingFee: SHIPPING_FEE,
			total: 230_000
		});
	});

	test('waives shipping exactly at the threshold', () => {
		expect(computeTotals([{ priceUnit: FREE_SHIPPING_THRESHOLD, quantity: 1 }])).toEqual({
			subtotal: FREE_SHIPPING_THRESHOLD,
			shippingFee: 0,
			total: FREE_SHIPPING_THRESHOLD
		});
	});

	test('still charges just below the threshold', () => {
		const totals = computeTotals([{ priceUnit: FREE_SHIPPING_THRESHOLD - 1, quantity: 1 }]);
		expect(totals.shippingFee).toBe(SHIPPING_FEE);
		expect(totals.total).toBe(FREE_SHIPPING_THRESHOLD - 1 + SHIPPING_FEE);
	});

	test('accumulates multiple lines', () => {
		expect(
			computeTotals([
				{ priceUnit: 300_000, quantity: 1 },
				{ priceUnit: 120_000, quantity: 2 }
			])
		).toEqual({ subtotal: 540_000, shippingFee: 0, total: 540_000 });
	});
});
