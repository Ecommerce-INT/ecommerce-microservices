import { sql } from "./db";
import type { FavouriteDto, FavouriteId } from "./types";

function formatDate(dateVal: any): string {
  if (dateVal instanceof Date) {
    return dateVal.toISOString().slice(0, 19);
  }
  return String(dateVal);
}

export class FavouriteService {
  async findAll(): Promise<FavouriteDto[]> {
    const rows = await sql`
      SELECT user_id, product_id, like_date
      FROM favourites
      ORDER BY like_date DESC
    `;

    return rows.map((r: any) => ({
      userId: r.user_id,
      productId: r.product_id,
      likeDate: formatDate(r.like_date),
    }));
  }

  async findById(id: FavouriteId): Promise<FavouriteDto | null> {
    const rows = await sql`
      SELECT user_id, product_id, like_date
      FROM favourites
      WHERE user_id = ${id.userId}
        AND product_id = ${id.productId}
        AND like_date = ${id.likeDate}::timestamp
      LIMIT 1
    `;

    if (rows.length === 0) {
      return null;
    }

    const r = rows[0];
    return {
      userId: r.user_id,
      productId: r.product_id,
      likeDate: formatDate(r.like_date),
    };
  }

  async save(dto: FavouriteDto): Promise<FavouriteDto> {
    await sql`
      INSERT INTO favourites (user_id, product_id, like_date)
      VALUES (${dto.userId}, ${dto.productId}, ${dto.likeDate}::timestamp)
      ON CONFLICT (user_id, product_id, like_date) DO UPDATE
      SET updated_at = CURRENT_TIMESTAMP
    `;

    return dto;
  }

  async update(dto: FavouriteDto): Promise<FavouriteDto> {
    return this.save(dto);
  }

  async deleteById(id: FavouriteId): Promise<boolean> {
    const result = await sql`
      DELETE FROM favourites
      WHERE user_id = ${id.userId}
        AND product_id = ${id.productId}
        AND like_date = ${id.likeDate}::timestamp
    `;

    return result.count >= 0;
  }
}
