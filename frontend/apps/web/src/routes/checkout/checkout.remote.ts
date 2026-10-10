import { error } from '@sveltejs/kit';
import { command, getRequestEvent } from '$app/server';
import { computeTotals } from '@ecommerce/lib/cart';
import { checkoutSchema } from '@ecommerce/lib/schemas';
import { createApi } from '@ecommerce/lib/server';
import { accessToken } from '@ecommerce/lib/server';

/**
 * Creates the order + payment for the client-side cart. Line prices are
 * re-read from the catalog so the charged fee is never client-supplied, and the
 * totals come from the shared `computeTotals` used by the cart UI.
 */
export const placeOrder = command(checkoutSchema, async (input) => {
	const { cookies } = getRequestEvent();
	const token = accessToken(cookies);
	if (!token) error(401, 'Your session has expired. Please sign in again.');

	const api = createApi({ token });
	const products = await Promise.all(input.items.map((line) => api.products.get(line.productId)));

	const totals = computeTotals(
		input.items.map((line, index) => ({
			priceUnit: products[index].priceUnit,
			quantity: line.quantity
		}))
	);

	// The order model carries a single product reference — the first line — and
	// the checkout total in `orderFee` (same contract the previous app used).
	const order = await api.orders.create({
		orderDate: new Date().toISOString(),
		orderDesc: input.shipping.note || 'EzBuy order',
		orderFee: totals.total,
		productId: products[0].productId
	});

	await api.payments.create({
		isPayed: input.paymentMethod !== 'COD',
		paymentStatus: input.paymentMethod === 'COD' ? 'NOT_STARTED' : 'IN_PROGRESS',
		orderId: order.orderId
	});

	return { orderId: order.orderId, total: totals.total };
});
