import { sql } from 'drizzle-orm';
import { bigint, bigserial, boolean, numeric, pgTable, text, timestamp, varchar } from 'drizzle-orm/pg-core';

export const notifications = pgTable('notifications', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  userId: varchar('user_id', { length: 255 }),
  title: varchar('title', { length: 255 }),
  message: text('message'),
  type: varchar('type', { length: 100 }),
  status: varchar('status', { length: 50 }).default('UNREAD'),
  createdAt: timestamp('created_at').default(sql`CURRENT_TIMESTAMP`)
});

export const paymentNotifications = pgTable('payment_notifications', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  paymentId: bigint('payment_id', { mode: 'number' }),
  userId: varchar('user_id', { length: 255 }),
  orderId: bigint('order_id', { mode: 'number' }),
  amount: numeric('amount'),
  isPayed: boolean('is_payed'),
  paymentStatus: varchar('payment_status', { length: 50 }),
  createdAt: timestamp('created_at').default(sql`CURRENT_TIMESTAMP`)
});

export type NotificationRow = typeof notifications.$inferSelect;
export type PaymentNotificationRow = typeof paymentNotifications.$inferSelect;
