import type { SessionUser } from '@ecommerce/lib/types';

declare global {
	namespace App {
		interface Locals {
			user: SessionUser | null;
		}
	}
}

export {};
