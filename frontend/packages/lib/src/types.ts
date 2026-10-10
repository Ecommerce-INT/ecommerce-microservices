export interface ApiResponse<T> {
	success: boolean;
	code: string;
	message?: string;
	data?: T;
	errors?: string[];
	path?: string;
	traceId?: string;
	timestamp: string;
}

export interface PaginatedResponse<T> {
	content: T[];
	totalElements: number;
	totalPages: number;
	size: number;
	number: number;
	first: boolean;
	last: boolean;
}

export interface TokenResponse {
	access_token: string;
	refresh_token: string;
	expires_in: number;
	refresh_expires_in: number;
	token_type: string;
	scope: string;
}

export interface Role {
	id: number;
	name: string;
}

export interface UserResponse {
	id: number;
	fullName: string;
	username: string;
	email: string;
	gender: string;
	phone?: string;
	avatar?: string;
	roles: Role[];
}

/** Normalized user shape the BFF keeps in its session cookie. */
export interface SessionUser {
	id: number;
	username: string;
	fullName?: string;
	email?: string;
	roles: string[];
}

export interface Category {
	categoryId: number;
	categoryTitle: string;
	imageUrl?: string;
	parentCategory?: Category;
	subCategories?: Category[];
	products?: Product[];
}

export interface Product {
	productId: number;
	productTitle: string;
	imageUrl?: string;
	sku?: string;
	priceUnit: number;
	quantity: number;
	category?: Category;
	description?: string;
}

export interface Cart {
	cartId: number;
	userId: number;
	orderItems?: OrderItem[];
}

export interface OrderItem {
	orderId: number;
	productId: number;
	quantity: number;
	unitPrice: number;
	product?: Product;
}

export interface Order {
	orderId: number;
	orderDate: string;
	orderDesc?: string;
	orderFee: number;
	productId: number;
	cart?: Cart;
}

export type PaymentStatus = 'IN_PROGRESS' | 'COMPLETED' | 'NOT_STARTED' | 'FAILED' | 'CANCELLED';

export interface Payment {
	paymentId: number;
	isPayed: boolean;
	paymentStatus: PaymentStatus;
	orderId: number;
}

export interface InventoryItem {
	id: number;
	skuCode: string;
	quantity: number;
}

export interface Favourite {
	userId: number;
	productId: number;
	product?: Product;
}

export interface Rating {
	id: number;
	productId: number;
	userId: string;
	ratingValue: number;
	comment?: string;
	createdAt: string;
}

export interface Notification {
	id: number;
	userId: number;
	message: string;
	isRead: boolean;
	createdAt: string;
}

/** A cart line held client-side (localStorage-persisted, never synced to the backend). */
export interface CartItem {
	product: Product;
	quantity: number;
}

export interface ProductPayload {
	productTitle: string;
	sku: string;
	imageUrl?: string | null;
	priceUnit: number;
	quantity: number;
	categoryId?: number | null;
	description?: string | null;
}

export interface CategoryPayload {
	categoryTitle: string;
	imageUrl?: string | null;
}

export interface CreateOrderPayload {
	orderDate: string;
	orderDesc?: string;
	orderFee: number;
	productId: number;
}

export interface CreatePaymentPayload {
	isPayed: boolean;
	paymentStatus: PaymentStatus;
	orderId: number;
}

export interface ProductListParams {
	page?: number;
	size?: number;
	sort?: string;
	categoryId?: number;
	search?: string;
}
