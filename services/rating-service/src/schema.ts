import { sql } from 'drizzle-orm';
import { bigint, bigserial, integer, pgTable, timestamp, varchar } from 'drizzle-orm/pg-core';

export const rating = pgTable('rating', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  content: varchar('content', { length: 255 }),
  ratingStar: integer('rating_star').notNull(),
  productId: bigint('product_id', { mode: 'number' }),
  productName: varchar('product_name', { length: 255 }),
  firstName: varchar('first_name', { length: 255 }),
  lastName: varchar('last_name', { length: 255 }),
  createdBy: varchar('created_by', { length: 255 }),
  createdOn: timestamp('created_on').default(sql`CURRENT_TIMESTAMP`),
  lastModifiedBy: varchar('last_modified_by', { length: 255 }),
  lastModifiedOn: timestamp('last_modified_on').default(sql`CURRENT_TIMESTAMP`)
});

export type RatingRow = typeof rating.$inferSelect;
export type NewRatingRow = typeof rating.$inferInsert;
