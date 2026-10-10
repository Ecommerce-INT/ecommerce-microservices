import { sql } from "./db";
import type { PromotionDetailVm, PromotionListVm, PromotionPostVm, PromotionPutVm } from "./types";

function mapRow(r: any): PromotionDetailVm {
  return {
    id: Number(r.id),
    name: r.name || "",
    slug: r.slug || "",
    description: r.description || "",
    couponCode: r.coupon_code || "",
    usageLimit: Number(r.usage_limit || 0),
    usageCount: Number(r.usage_count || 0),
    discountType: r.discount_type || "PERCENTAGE",
    applyTo: r.apply_to || "ALL",
    usageType: r.usage_type || "UNLIMITED",
    discountPercentage: Number(r.discount_percentage || 0),
    discountAmount: Number(r.discount_amount || 0),
    isActive: Boolean(r.is_active),
    startDate: r.start_date instanceof Date ? r.start_date.toISOString() : (r.start_date || null),
    endDate: r.end_date instanceof Date ? r.end_date.toISOString() : (r.end_date || null),
    brands: [],
    categories: [],
    products: [],
  };
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^\w ]+/g, "")
    .replace(/ +/g, "-");
}

export class PromotionService {
  async getPromotions(
    pageNo: number,
    pageSize: number,
    promotionName = "",
    couponCode = "",
    startDate?: string,
    endDate?: string
  ): Promise<PromotionListVm> {
    const offset = pageNo * pageSize;
    const namePattern = `%${promotionName.toLowerCase()}%`;
    const codePattern = `%${couponCode.toLowerCase()}%`;

    const countRes = await sql`
      SELECT COUNT(*) as count FROM promotion
      WHERE LOWER(COALESCE(name, '')) LIKE ${namePattern}
        AND LOWER(COALESCE(coupon_code, '')) LIKE ${codePattern}
    `;
    const totalElements = Number(countRes[0]?.count || 0);
    const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

    const rows = await sql`
      SELECT * FROM promotion
      WHERE LOWER(COALESCE(name, '')) LIKE ${namePattern}
        AND LOWER(COALESCE(coupon_code, '')) LIKE ${codePattern}
      ORDER BY id DESC
      LIMIT ${pageSize} OFFSET ${offset}
    `;

    return {
      promotionDetailVmList: rows.map(mapRow),
      totalElements,
      totalPages,
    };
  }

  async getPromotion(id: number): Promise<PromotionDetailVm | null> {
    const rows = await sql`SELECT * FROM promotion WHERE id = ${id}`;
    if (rows.length === 0) return null;
    return mapRow(rows[0]);
  }

  async createPromotion(req: PromotionPostVm): Promise<PromotionDetailVm> {
    const slug = req.slug || slugify(req.name);
    const rows = await sql`
      INSERT INTO promotion (
        name, slug, description, coupon_code, usage_limit,
        discount_type, apply_to, usage_type, discount_percentage, discount_amount,
        is_active, start_date, end_date, created_on, last_modified_on
      ) VALUES (
        ${req.name}, ${slug}, ${req.description}, ${req.couponCode}, ${req.usageLimit},
        ${req.discountType}, ${req.applyTo}, ${req.usageType}, ${req.discountPercentage}, ${req.discountAmount},
        ${req.isActive}, ${req.startDate ? new Date(req.startDate) : null}, ${req.endDate ? new Date(req.endDate) : null},
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
      )
      RETURNING *
    `;

    return mapRow(rows[0]);
  }

  async updatePromotion(req: PromotionPutVm): Promise<PromotionDetailVm> {
    const slug = req.slug || slugify(req.name);
    const rows = await sql`
      UPDATE promotion
      SET name = ${req.name},
          slug = ${slug},
          description = ${req.description},
          coupon_code = ${req.couponCode},
          usage_limit = ${req.usageLimit},
          discount_type = ${req.discountType},
          apply_to = ${req.applyTo},
          usage_type = ${req.usageType},
          discount_percentage = ${req.discountPercentage},
          discount_amount = ${req.discountAmount},
          is_active = ${req.isActive},
          start_date = ${req.startDate ? new Date(req.startDate) : null},
          end_date = ${req.endDate ? new Date(req.endDate) : null},
          last_modified_on = CURRENT_TIMESTAMP
      WHERE id = ${req.id}
      RETURNING *
    `;

    return mapRow(rows[0]);
  }

  async deletePromotion(id: number): Promise<void> {
    await sql`DELETE FROM promotion WHERE id = ${id}`;
  }
}
