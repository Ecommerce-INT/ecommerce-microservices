import { and, count, desc, eq, sql } from 'drizzle-orm';
import { db } from './db';
import { promotion } from './schema';
import type { PromotionDetailVm, PromotionListVm, PromotionPostVm, PromotionPutVm } from './types';

function mapRow(r: typeof promotion.$inferSelect): PromotionDetailVm {
  return {
    id: r.id,
    name: r.name || '',
    slug: r.slug || '',
    description: r.description || '',
    couponCode: r.couponCode || '',
    usageLimit: Number(r.usageLimit || 0),
    usageCount: Number(r.usageCount || 0),
    discountType: r.discountType || 'PERCENTAGE',
    applyTo: r.applyTo || 'ALL',
    usageType: r.usageType || 'UNLIMITED',
    discountPercentage: Number(r.discountPercentage || 0),
    discountAmount: Number(r.discountAmount || 0),
    isActive: Boolean(r.isActive),
    startDate:
      r.startDate instanceof Date ? r.startDate.toISOString() : r.startDate || null,
    endDate: r.endDate instanceof Date ? r.endDate.toISOString() : r.endDate || null,
    brands: [],
    categories: [],
    products: []
  };
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^\w ]+/g, '')
    .replace(/ +/g, '-');
}

export class PromotionService {
  async getPromotions(
    pageNo: number,
    pageSize: number,
    promotionName = '',
    couponCode = '',
    startDate?: string,
    endDate?: string
  ): Promise<PromotionListVm> {
    try {
      const offset = pageNo * pageSize;
      const namePattern = `%${promotionName.toLowerCase()}%`;
      const codePattern = `%${couponCode.toLowerCase()}%`;

      const filter = and(
        sql`LOWER(COALESCE(${promotion.name}, '')) LIKE ${namePattern}`,
        sql`LOWER(COALESCE(${promotion.couponCode}, '')) LIKE ${codePattern}`
      );

      const [countRow] = await db.select({ total: count() }).from(promotion).where(filter);
      const totalElements = Number(countRow?.total ?? 0);
      const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

      const rows = await db
        .select()
        .from(promotion)
        .where(filter)
        .orderBy(desc(promotion.id))
        .limit(pageSize)
        .offset(offset);

      return {
        promotionDetailVmList: rows.map(mapRow),
        totalElements,
        totalPages
      };
    } catch {
      return {
        promotionDetailVmList: [],
        totalElements: 0,
        totalPages: 0
      };
    }
  }

  async getPromotion(id: number): Promise<PromotionDetailVm | null> {
    const rows = await db.select().from(promotion).where(eq(promotion.id, id));
    if (rows.length === 0) return null;
    return mapRow(rows[0]);
  }

  async createPromotion(req: PromotionPostVm): Promise<PromotionDetailVm> {
    const slug = req.slug || slugify(req.name);
    const [row] = await db
      .insert(promotion)
      .values({
        name: req.name,
        slug,
        description: req.description,
        couponCode: req.couponCode,
        usageLimit: req.usageLimit,
        discountType: req.discountType,
        applyTo: req.applyTo,
        usageType: req.usageType,
        discountPercentage: req.discountPercentage,
        discountAmount: req.discountAmount,
        isActive: req.isActive,
        startDate: req.startDate ? new Date(req.startDate) : null,
        endDate: req.endDate ? new Date(req.endDate) : null
      })
      .returning();

    return mapRow(row);
  }

  async updatePromotion(req: PromotionPutVm): Promise<PromotionDetailVm> {
    const slug = req.slug || slugify(req.name);
    const [row] = await db
      .update(promotion)
      .set({
        name: req.name,
        slug,
        description: req.description,
        couponCode: req.couponCode,
        usageLimit: req.usageLimit,
        discountType: req.discountType,
        applyTo: req.applyTo,
        usageType: req.usageType,
        discountPercentage: req.discountPercentage,
        discountAmount: req.discountAmount,
        isActive: req.isActive,
        startDate: req.startDate ? new Date(req.startDate) : null,
        endDate: req.endDate ? new Date(req.endDate) : null,
        lastModifiedOn: sql`CURRENT_TIMESTAMP`
      })
      .where(eq(promotion.id, req.id))
      .returning();

    return mapRow(row);
  }

  async deletePromotion(id: number): Promise<void> {
    await db.delete(promotion).where(eq(promotion.id, id));
  }
}
