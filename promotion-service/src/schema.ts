import { sql } from 'drizzle-orm';
import { bigint, bigserial, boolean, integer, numeric, pgTable, timestamp, varchar } from 'drizzle-orm/pg-core';

export const promotion = pgTable('promotion', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  name: varchar('name', { length: 255 }).notNull(),
  slug: varchar('slug', { length: 255 }).notNull(),
  description: varchar('description', { length: 255 }),
  couponCode: varchar('coupon_code', { length: 255 }),
  discountPercentage: bigint('discount_percentage', { mode: 'number' }).notNull().default(0),
  discountAmount: bigint('discount_amount', { mode: 'number' }).notNull().default(0),
  isActive: boolean('is_active').notNull().default(true),
  startDate: timestamp('start_date', { withTimezone: true }),
  endDate: timestamp('end_date', { withTimezone: true }),
  discountType: varchar('discount_type', { length: 25 }).notNull().default('PERCENTAGE'),
  usageLimit: integer('usage_limit').default(100),
  usageCount: integer('usage_count').notNull().default(0),
  usageType: varchar('usage_type', { length: 25 }).notNull().default('UNLIMITED'),
  applyTo: varchar('apply_to', { length: 25 }).notNull().default('ALL'),
  minimumOrderPurchaseAmount: bigint('minimum_order_purchase_amount', { mode: 'number' }).default(0),
  createdBy: varchar('created_by', { length: 255 }),
  createdOn: timestamp('created_on', { withTimezone: true }).default(sql`CURRENT_TIMESTAMP`),
  lastModifiedBy: varchar('last_modified_by', { length: 255 }),
  lastModifiedOn: timestamp('last_modified_on', { withTimezone: true }).default(sql`CURRENT_TIMESTAMP`)
});

export const taxClass = pgTable('tax_class', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  name: varchar('name', { length: 255 }).notNull().unique()
});

export const taxRate = pgTable('tax_rate', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  rate: numeric('rate', { precision: 10, scale: 4 }).notNull(),
  countryId: bigint('country_id', { mode: 'number' }).notNull(),
  stateOrProvinceId: bigint('state_or_province_id', { mode: 'number' }),
  zipCode: varchar('zip_code', { length: 50 }),
  taxClassId: bigint('tax_class_id', { mode: 'number' }).references(() => taxClass.id)
});

export type PromotionRow = typeof promotion.$inferSelect;
export type TaxClassRow = typeof taxClass.$inferSelect;
export type TaxRateRow = typeof taxRate.$inferSelect;
