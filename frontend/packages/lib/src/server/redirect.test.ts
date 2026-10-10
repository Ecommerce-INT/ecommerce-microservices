import { describe, expect, test } from 'bun:test';
import { safeRedirect } from './redirect';

describe('safeRedirect', () => {
	test('accepts in-app absolute paths', () => {
		expect(safeRedirect('/orders')).toBe('/orders');
		expect(safeRedirect('/orders?page=2')).toBe('/orders?page=2');
	});

	test('rejects absolute URLs and protocol-relative paths', () => {
		expect(safeRedirect('https://evil.example')).toBe('/');
		expect(safeRedirect('//evil.example')).toBe('/');
		expect(safeRedirect('orders')).toBe('/');
	});

	test('falls back to the root for empty values', () => {
		expect(safeRedirect(null)).toBe('/');
		expect(safeRedirect(undefined)).toBe('/');
		expect(safeRedirect('')).toBe('/');
	});
});
