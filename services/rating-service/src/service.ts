import { and, count, desc, eq, gte, lte, sql } from 'drizzle-orm';
import { db } from './db';
import { rating } from './schema';
import type { RatingListVm, RatingPostVm, RatingVm, ResponeStatusVm } from './types';

function mapRow(r: typeof rating.$inferSelect): RatingVm {
  return {
    id: r.id,
    content: r.content || '',
    star: r.ratingStar,
    productId: r.productId ?? 0,
    productName: r.productName || '',
    createdBy: r.createdBy || '',
    lastName: r.lastName || '',
    firstName: r.firstName || '',
    createdOn: r.createdOn instanceof Date ? r.createdOn.toISOString() : String(r.createdOn || '')
  };
}

export class RatingService {
  async getRatingListByProductId(
    productId: number,
    pageNo: number,
    pageSize: number
  ): Promise<RatingListVm> {
    const offset = pageNo * pageSize;

    const [countRow] = await db
      .select({ total: count() })
      .from(rating)
      .where(eq(rating.productId, productId));
    const totalElements = Number(countRow?.total ?? 0);
    const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

    const rows = await db
      .select()
      .from(rating)
      .where(eq(rating.productId, productId))
      .orderBy(desc(rating.createdOn))
      .limit(pageSize)
      .offset(offset);

    return {
      ratingList: rows.map(mapRow),
      totalElements,
      totalPages
    };
  }

  async getRatingListWithFilter(
    proName: string,
    cusName: string,
    message: string,
    createdFrom: string,
    createdTo: string,
    pageNo: number,
    pageSize: number
  ): Promise<RatingListVm> {
    const offset = pageNo * pageSize;
    const proPattern = `%${(proName || '').toLowerCase()}%`;
    const cusPattern = `%${(cusName || '').toLowerCase()}%`;
    const msgPattern = `%${(message || '').toLowerCase()}%`;

    const fromDate = createdFrom ? new Date(createdFrom) : new Date(0);
    const toDate = createdTo ? new Date(createdTo) : new Date();

    const filter = and(
      sql`LOWER(COALESCE(${rating.productName}, '')) LIKE ${proPattern}`,
      sql`LOWER(CONCAT(COALESCE(${rating.firstName}, ''), ' ', COALESCE(${rating.lastName}, ''))) LIKE ${cusPattern}`,
      sql`LOWER(COALESCE(${rating.content}, '')) LIKE ${msgPattern}`,
      gte(rating.createdOn, fromDate),
      lte(rating.createdOn, toDate)
    );

    const [countRow] = await db.select({ total: count() }).from(rating).where(filter);
    const totalElements = Number(countRow?.total ?? 0);
    const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

    const rows = await db
      .select()
      .from(rating)
      .where(filter)
      .orderBy(desc(rating.createdOn))
      .limit(pageSize)
      .offset(offset);

    return {
      ratingList: rows.map(mapRow),
      totalElements,
      totalPages
    };
  }

  async createRating(
    req: RatingPostVm,
    userClaims: { userId?: string; firstName?: string; lastName?: string }
  ): Promise<RatingVm> {
    const userId = userClaims.userId || 'anonymous';
    const firstName = userClaims.firstName || 'Customer';
    const lastName = userClaims.lastName || '';

    const [row] = await db
      .insert(rating)
      .values({
        content: req.content,
        ratingStar: req.star,
        productId: req.productId,
        productName: req.productName,
        firstName,
        lastName,
        createdBy: userId
      })
      .returning();

    return mapRow(row);
  }

  async deleteRating(id: number): Promise<ResponeStatusVm> {
    await db.delete(rating).where(eq(rating.id, id));
    return {
      title: 'Delete Rating',
      message: 'The request has been processed successfully',
      statusCode: '200'
    };
  }

  async calculateAverageStar(productId: number): Promise<number> {
    const [row] = await db
      .select({
        totalStars: sql<number>`COALESCE(SUM(${rating.ratingStar}), 0)`,
        totalRatings: count()
      })
      .from(rating)
      .where(eq(rating.productId, productId));

    const totalStars = Number(row?.totalStars ?? 0);
    const totalRatings = Number(row?.totalRatings ?? 0);

    if (totalRatings === 0) {
      return 0.0;
    }

    return Number((totalStars / totalRatings).toFixed(2));
  }
}
