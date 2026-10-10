import { describe, expect, test } from 'bun:test';
import { formatDate, formatPrice, getInitials, truncateText } from './format';

describe('formatPrice', () => {
	test('formats VND with vi-VN grouping and symbol', () => {
		expect(formatPrice(1234567)).toMatch(/^1\.234\.567\s₫$/);
		expect(formatPrice(0)).toMatch(/^0\s₫$/);
	});

	test('formats non-VND currencies with en-US grouping', () => {
		expect(formatPrice(1234567, 'USD')).toBe('$1,234,567.00');
	});
});

describe('formatDate', () => {
	test('formats as HH:mm dd/MM/yyyy in local time', () => {
		expect(formatDate('2026-03-05T21:30:00')).toBe('21:30 05/03/2026');
	});
});

describe('truncateText', () => {
	test('returns text at or below the limit unchanged', () => {
		expect(truncateText('hello', 5)).toBe('hello');
		expect(truncateText('hey', 5)).toBe('hey');
	});

	test('cuts longer text and appends an ellipsis', () => {
		expect(truncateText('hello world', 5)).toBe('hello...');
	});
});

describe('getInitials', () => {
	test('takes the first letters of the first two words, uppercased', () => {
		expect(getInitials('Nguyen Van An')).toBe('NV');
		expect(getInitials('nguyen van an')).toBe('NV');
	});

	test('handles a single name', () => {
		expect(getInitials('An')).toBe('A');
	});
});
