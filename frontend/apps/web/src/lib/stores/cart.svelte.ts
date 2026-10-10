import { computeTotals } from '@ecommerce/lib/cart';
import type { CartItem, Product } from '@ecommerce/lib/types';

const STORAGE_KEY = 'ecommerce-cart';

function readStoredItems(): CartItem[] {
	if (typeof localStorage === 'undefined') return [];
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return [];
		const parsed: unknown = JSON.parse(raw);
		// Tolerates the old zustand-persist shape ({ state: { items } }), a
		// plain { items } object, or a bare array.
		if (Array.isArray(parsed)) return parsed as CartItem[];
		const candidate =
			parsed && typeof parsed === 'object' && 'state' in parsed
				? (parsed as { state?: { items?: CartItem[] } }).state
				: (parsed as { items?: CartItem[] });
		return Array.isArray(candidate?.items) ? candidate.items : [];
	} catch {
		return [];
	}
}

class CartStore {
	items = $state<CartItem[]>([]);
	hydrated = $state(false);

	totalItems = $derived(this.items.reduce((sum, item) => sum + item.quantity, 0));
	totals = $derived(
		computeTotals(
			this.items.map((item) => ({ priceUnit: item.product.priceUnit, quantity: item.quantity }))
		)
	);
	totalPrice = $derived(this.totals.subtotal);
	shippingFee = $derived(this.totals.shippingFee);
	grandTotal = $derived(this.totals.total);

	/** Reads localStorage on the client. Called once from the root layout on mount. */
	hydrate() {
		if (this.hydrated) return;
		this.items = readStoredItems();
		this.hydrated = true;
	}

	add(product: Product, quantity = 1) {
		const existing = this.items.find((item) => item.product.productId === product.productId);
		if (existing) {
			existing.quantity += quantity;
		} else {
			this.items.push({ product, quantity });
		}
		this.persist();
	}

	remove(productId: number) {
		this.items = this.items.filter((item) => item.product.productId !== productId);
		this.persist();
	}

	updateQuantity(productId: number, quantity: number) {
		if (quantity <= 0) {
			this.remove(productId);
			return;
		}
		const item = this.items.find((entry) => entry.product.productId === productId);
		if (item) {
			item.quantity = quantity;
			this.persist();
		}
	}

	clear() {
		this.items = [];
		this.persist();
	}

	private persist() {
		if (typeof localStorage === 'undefined') return;
		try {
			localStorage.setItem(
				STORAGE_KEY,
				JSON.stringify({ state: { items: this.items }, version: 0 })
			);
		} catch {
			// Storage full / disabled — the cart simply stays in memory.
		}
	}
}

export const cart = new CartStore();
