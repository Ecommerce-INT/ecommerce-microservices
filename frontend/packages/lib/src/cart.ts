/** Subtotal at or above this amount ships free. */
export const FREE_SHIPPING_THRESHOLD = 500_000;

/** Flat fee applied below the free-shipping threshold. */
export const SHIPPING_FEE = 30_000;

export interface CartLine {
	priceUnit: number;
	quantity: number;
}

export interface CartTotals {
	subtotal: number;
	shippingFee: number;
	total: number;
}

/**
 * Single definition of cart money math, shared by the client cart store and the
 * server-side checkout command so both always agree.
 */
export function computeTotals(lines: CartLine[]): CartTotals {
	const subtotal = lines.reduce((sum, line) => sum + line.priceUnit * line.quantity, 0);
	const shippingFee = lines.length === 0 || subtotal >= FREE_SHIPPING_THRESHOLD ? 0 : SHIPPING_FEE;
	return { subtotal, shippingFee, total: subtotal + shippingFee };
}
