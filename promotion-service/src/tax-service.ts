import { and, asc, count, desc, eq, inArray, sql } from 'drizzle-orm';
import { db } from './db';
import { taxClass, taxRate } from './schema';
import type {
  TaxClassPostVm,
  TaxClassVm,
  TaxClassListVm,
  TaxRatePostVm,
  TaxRateVm,
  TaxRateListVm
} from './tax-types';

function mapTaxClass(r: typeof taxClass.$inferSelect): TaxClassVm {
  return {
    id: r.id,
    name: r.name
  };
}

function mapTaxRate(r: {
  tax_rate: typeof taxRate.$inferSelect;
  tax_class: typeof taxClass.$inferSelect | null;
}): TaxRateVm {
  const row = r.tax_rate;
  return {
    id: row.id,
    rate: Number(row.rate),
    countryId: row.countryId,
    stateOrProvinceId: row.stateOrProvinceId ?? null,
    zipCode: row.zipCode,
    taxClassId: row.taxClassId ?? 0,
    taxClassName: r.tax_class?.name ?? null
  };
}

const taxRateWithClass = {
  tax_rate: taxRate,
  tax_class: taxClass
};

export class TaxService {
  // ================= Tax Classes =================

  async findAllTaxClasses(): Promise<TaxClassVm[]> {
    const rows = await db.select().from(taxClass).orderBy(asc(taxClass.id));
    return rows.map(mapTaxClass);
  }

  async findTaxClassById(id: number): Promise<TaxClassVm | null> {
    const rows = await db.select().from(taxClass).where(eq(taxClass.id, id)).limit(1);
    if (rows.length === 0) return null;
    return mapTaxClass(rows[0]);
  }

  async createTaxClass(vm: TaxClassPostVm): Promise<TaxClassVm> {
    const [row] = await db.insert(taxClass).values({ name: vm.name }).returning();
    return mapTaxClass(row);
  }

  async updateTaxClass(id: number, vm: TaxClassPostVm): Promise<boolean> {
    const rows = await db
      .update(taxClass)
      .set({ name: vm.name })
      .where(eq(taxClass.id, id))
      .returning({ id: taxClass.id });
    return rows.length > 0;
  }

  async deleteTaxClass(id: number): Promise<boolean> {
    const rows = await db.delete(taxClass).where(eq(taxClass.id, id)).returning({ id: taxClass.id });
    return rows.length > 0;
  }

  async getPageableTaxClasses(pageNo = 0, pageSize = 10): Promise<TaxClassListVm> {
    try {
      const offset = pageNo * pageSize;
      const [countRow] = await db.select({ total: count() }).from(taxClass);
      const total = Number(countRow?.total ?? 0);

      const rows = await db
        .select()
        .from(taxClass)
        .orderBy(asc(taxClass.id))
        .limit(pageSize)
        .offset(offset);

      const totalPages = Math.ceil(total / pageSize);
      return {
        taxClasses: rows.map(mapTaxClass),
        pageNo,
        pageSize,
        totalElements: total,
        totalPages,
        isLast: pageNo >= totalPages - 1
      };
    } catch {
      return {
        taxClasses: [],
        pageNo,
        pageSize,
        totalElements: 0,
        totalPages: 0,
        isLast: true
      };
    }
  }

  // ================= Tax Rates =================

  async findAllTaxRates(): Promise<TaxRateVm[]> {
    const rows = await db
      .select(taxRateWithClass)
      .from(taxRate)
      .leftJoin(taxClass, eq(taxRate.taxClassId, taxClass.id))
      .orderBy(asc(taxRate.id));
    return rows.map(mapTaxRate);
  }

  async findTaxRateById(id: number): Promise<TaxRateVm | null> {
    const rows = await db
      .select(taxRateWithClass)
      .from(taxRate)
      .leftJoin(taxClass, eq(taxRate.taxClassId, taxClass.id))
      .where(eq(taxRate.id, id))
      .limit(1);
    if (rows.length === 0) return null;
    return mapTaxRate(rows[0]);
  }

  async createTaxRate(vm: TaxRatePostVm): Promise<TaxRateVm> {
    const [row] = await db
      .insert(taxRate)
      .values({
        rate: String(vm.rate),
        countryId: vm.countryId,
        stateOrProvinceId: vm.stateOrProvinceId ?? null,
        zipCode: vm.zipCode ?? null,
        taxClassId: vm.taxClassId
      })
      .returning();
    return {
      id: row.id,
      rate: Number(row.rate),
      countryId: row.countryId,
      stateOrProvinceId: row.stateOrProvinceId ?? null,
      zipCode: row.zipCode,
      taxClassId: row.taxClassId ?? 0
    };
  }

  async updateTaxRate(id: number, vm: TaxRatePostVm): Promise<boolean> {
    const rows = await db
      .update(taxRate)
      .set({
        rate: String(vm.rate),
        countryId: vm.countryId,
        stateOrProvinceId: vm.stateOrProvinceId ?? null,
        zipCode: vm.zipCode ?? null,
        taxClassId: vm.taxClassId
      })
      .where(eq(taxRate.id, id))
      .returning({ id: taxRate.id });
    return rows.length > 0;
  }

  async deleteTaxRate(id: number): Promise<boolean> {
    const rows = await db.delete(taxRate).where(eq(taxRate.id, id)).returning({ id: taxRate.id });
    return rows.length > 0;
  }

  async getPageableTaxRates(pageNo = 0, pageSize = 10): Promise<TaxRateListVm> {
    const offset = pageNo * pageSize;
    const [countRow] = await db.select({ total: count() }).from(taxRate);
    const total = Number(countRow?.total ?? 0);

    const rows = await db
      .select(taxRateWithClass)
      .from(taxRate)
      .leftJoin(taxClass, eq(taxRate.taxClassId, taxClass.id))
      .orderBy(asc(taxRate.id))
      .limit(pageSize)
      .offset(offset);

    const totalPages = Math.ceil(total / pageSize);
    return {
      taxRates: rows.map(mapTaxRate),
      pageNo,
      pageSize,
      totalElements: total,
      totalPages,
      isLast: pageNo >= totalPages - 1
    };
  }

  async getTaxPercent(
    taxClassId: number,
    countryId: number,
    stateOrProvinceId?: number | null,
    zipCode?: string | null
  ): Promise<number> {
    const rows = await db
      .select({ rate: taxRate.rate })
      .from(taxRate)
      .where(and(eq(taxRate.taxClassId, taxClassId), eq(taxRate.countryId, countryId)))
      .orderBy(desc(taxRate.id))
      .limit(1);
    if (rows.length === 0) return 0.0;
    return Number(rows[0].rate);
  }

  async getBatchTaxRates(
    taxClassIds: number[],
    countryId: number,
    stateOrProvinceId?: number | null,
    zipCode?: string | null
  ): Promise<TaxRateVm[]> {
    if (taxClassIds.length === 0) return [];

    const rows = await db
      .select(taxRateWithClass)
      .from(taxRate)
      .leftJoin(taxClass, eq(taxRate.taxClassId, taxClass.id))
      .where(
        and(
          eq(taxRate.countryId, countryId),
          inArray(taxRate.taxClassId, taxClassIds)
        )
      );

    return rows.map(mapTaxRate);
  }
}
