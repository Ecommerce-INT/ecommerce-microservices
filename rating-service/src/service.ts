import { sql } from "./db";
import type { RatingListVm, RatingPostVm, RatingVm, ResponeStatusVm } from "./types";

function mapRow(r: any): RatingVm {
  return {
    id: Number(r.id),
    content: r.content || "",
    star: Number(r.rating_star),
    productId: Number(r.product_id),
    productName: r.product_name || "",
    createdBy: r.created_by || "",
    lastName: r.last_name || "",
    firstName: r.first_name || "",
    createdOn: r.created_on instanceof Date ? r.created_on.toISOString() : String(r.created_on || ""),
  };
}

export class RatingService {
  async getRatingListByProductId(productId: number, pageNo: number, pageSize: number): Promise<RatingListVm> {
    const offset = pageNo * pageSize;

    const countRes = await sql`
      SELECT COUNT(*) as count FROM rating WHERE product_id = ${productId}
    `;
    const totalElements = Number(countRes[0]?.count || 0);
    const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

    const rows = await sql`
      SELECT id, content, rating_star, product_id, product_name, first_name, last_name, created_by, created_on
      FROM rating
      WHERE product_id = ${productId}
      ORDER BY created_on DESC
      LIMIT ${pageSize} OFFSET ${offset}
    `;

    return {
      ratingList: rows.map(mapRow),
      totalElements,
      totalPages,
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
    const proPattern = `%${(proName || "").toLowerCase()}%`;
    const cusPattern = `%${(cusName || "").toLowerCase()}%`;
    const msgPattern = `%${(message || "").toLowerCase()}%`;

    const fromDate = createdFrom ? new Date(createdFrom) : new Date(0);
    const toDate = createdTo ? new Date(createdTo) : new Date();

    const countRes = await sql`
      SELECT COUNT(*) as count FROM rating
      WHERE LOWER(COALESCE(product_name, '')) LIKE ${proPattern}
        AND LOWER(CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, ''))) LIKE ${cusPattern}
        AND LOWER(COALESCE(content, '')) LIKE ${msgPattern}
        AND created_on BETWEEN ${fromDate} AND ${toDate}
    `;
    const totalElements = Number(countRes[0]?.count || 0);
    const totalPages = pageSize > 0 ? Math.ceil(totalElements / pageSize) : 0;

    const rows = await sql`
      SELECT id, content, rating_star, product_id, product_name, first_name, last_name, created_by, created_on
      FROM rating
      WHERE LOWER(COALESCE(product_name, '')) LIKE ${proPattern}
        AND LOWER(CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, ''))) LIKE ${cusPattern}
        AND LOWER(COALESCE(content, '')) LIKE ${msgPattern}
        AND created_on BETWEEN ${fromDate} AND ${toDate}
      ORDER BY created_on DESC
      LIMIT ${pageSize} OFFSET ${offset}
    `;

    return {
      ratingList: rows.map(mapRow),
      totalElements,
      totalPages,
    };
  }

  async createRating(
    req: RatingPostVm,
    userClaims: { userId?: string; firstName?: string; lastName?: string }
  ): Promise<RatingVm> {
    const userId = userClaims.userId || "anonymous";
    const firstName = userClaims.firstName || "Customer";
    const lastName = userClaims.lastName || "";

    const rows = await sql`
      INSERT INTO rating (
        content, rating_star, product_id, product_name,
        first_name, last_name, created_by, created_on
      ) VALUES (
        ${req.content}, ${req.star}, ${req.productId}, ${req.productName},
        ${firstName}, ${lastName}, ${userId}, CURRENT_TIMESTAMP
      )
      RETURNING id, content, rating_star, product_id, product_name, first_name, last_name, created_by, created_on
    `;

    return mapRow(rows[0]);
  }

  async deleteRating(id: number): Promise<ResponeStatusVm> {
    await sql`
      DELETE FROM rating WHERE id = ${id}
    `;
    return {
      title: "Delete Rating",
      message: "The request has been processed successfully",
      statusCode: "200",
    };
  }

  async calculateAverageStar(productId: number): Promise<number> {
    const res = await sql`
      SELECT COALESCE(SUM(rating_star), 0) as total_stars, COUNT(*) as total_ratings
      FROM rating
      WHERE product_id = ${productId}
    `;

    const totalStars = Number(res[0]?.total_stars || 0);
    const totalRatings = Number(res[0]?.total_ratings || 0);

    if (totalRatings === 0) {
      return 0.0;
    }

    return Number((totalStars / totalRatings).toFixed(2));
  }
}
