import { sql } from 'drizzle-orm';
import { bigint, bigserial, pgTable, text, varchar } from 'drizzle-orm/pg-core';

export const media = pgTable('media', {
  id: bigserial('id', { mode: 'number' }).primaryKey(),
  caption: varchar('caption', { length: 255 }),
  fileName: varchar('file_name', { length: 255 }),
  filePath: text('file_path'),
  mediaType: varchar('media_type', { length: 128 })
});

export type MediaRow = typeof media.$inferSelect;
