import { sql } from "./db";
import type {
  TaxClassPostVm,
  TaxClassVm,
  TaxClassListVm,
  TaxRatePostVm,
  TaxRateVm,
  TaxRateListVm,
} from "./tax-types";

export class TaxService {
  // ================= Tax Classes =================

  async findAllTaxClasses(): Promise<TaxClassVm[]> {
    const rows = await sql`
      SELECT id, name FROM tax_class ORDER BY id ASC
    `;
    return rows.map((r: any) => ({
      id: Number(r.id),
      name: r.name,
    }));
  }

  async findTaxClassById(id: number): Promise<TaxClassVm | null> {
    const rows = await sql`
      SELECT id, name FROM tax_class WHERE id = ${id} LIMIT 1
    `;
    if (rows.length === 0) return null;
    return {
      id: Number(rows[0].id),
      name: rows[0].name,
    };
  }

  async createTaxClass(vm: TaxClassPostVm): Promise<TaxClassVm> {
    const rows = await sql`
      INSERT INTO tax_class (name)
      VALUES (${vm.name})
      RETURNING id, name
    `;
    return {
      id: Number(rows[0].id),
      name: rows[0].name,
    };
  }

  async updateTaxClass(id: number, vm: TaxClassPostVm): Promise<boolean> {
    const res = await sql`
      UPDATE tax_class
      SET name = ${vm.name}
      WHERE id = ${id}
    `;
    return res.count > 0;
  }

  async deleteTaxClass(id: number): Promise<boolean> {
    const res = await sql`
      DELETE FROM tax_class WHERE id = ${id}
    `;
    return res.count > 0;
  }

  async getPageableTaxClasses(pageNo = 0, pageSize = 10): Promise<TaxClassListVm> {
    const offset = pageNo * pageSize;
    const [countRow] = await sql`SELECT COUNT(*)::int as total FROM tax_class`;
    const total = countRow.total || 0;

    const rows = await sql`
      SELECT id, name
      FROM tax_class
      ORDER BY id ASC
      LIMIT ${pageSize} OFFSET ${offset}
    `;

    const totalPages = Math.ceil(total / pageSize);
    return {
      taxClasses: rows.map((r: any) => ({ id: Number(r.id), name: r.name })),
      pageNo,
      pageSize,
      totalElements: total,
      totalPages,
      isLast: pageNo >= totalPages - 1,
    };
  }

  // ================= Tax Rates =================

  async findAllTaxRates(): Promise<TaxRateVm[]> {
    const rows = await sql`
      SELECT tr.id, tr.rate, tr.country_id, tr.state_or_province_id, tr.zip_code, tr.tax_class_id, tc.name as tax_class_name
      FROM tax_rate tr
      LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
      ORDER BY tr.id ASC
    `;
    return rows.map((r: any) => ({
      id: Number(r.id),
      rate: Number(r.rate),
      countryId: Number(r.country_id),
      stateOrProvinceId: r.state_or_province_id ? Number(r.state_or_province_id) : null,
      zipCode: r.zip_code,
      taxClassId: Number(r.tax_class_id),
      taxClassName: r.tax_class_name,
    }));
  }

  async findTaxRateById(id: number): Promise<TaxRateVm | null> {
    const rows = await sql`
      SELECT tr.id, tr.rate, tr.country_id, tr.state_or_province_id, tr.zip_code, tr.tax_class_id, tc.name as tax_class_name
      FROM tax_rate tr
      LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
      WHERE tr.id = ${id}
      LIMIT 1
    `;
    if (rows.length === 0) return null;
    const r = rows[0];
    return {
      id: Number(r.id),
      rate: Number(r.rate),
      countryId: Number(r.country_id),
      stateOrProvinceId: r.state_or_province_id ? Number(r.state_or_province_id) : null,
      zipCode: r.zip_code,
      taxClassId: Number(r.tax_class_id),
      taxClassName: r.tax_class_name,
    };
  }

  async createTaxRate(vm: TaxRatePostVm): Promise<TaxRateVm> {
    const rows = await sql`
      INSERT INTO tax_rate (rate, country_id, state_or_province_id, zip_code, tax_class_id)
      VALUES (${vm.rate}, ${vm.countryId}, ${vm.stateOrProvinceId ?? null}, ${vm.zipCode ?? null}, ${vm.taxClassId})
      RETURNING id, rate, country_id, state_or_province_id, zip_code, tax_class_id
    `;
    const r = rows[0];
    return {
      id: Number(r.id),
      rate: Number(r.rate),
      countryId: Number(r.country_id),
      stateOrProvinceId: r.state_or_province_id ? Number(r.state_or_province_id) : null,
      zipCode: r.zip_code,
      taxClassId: Number(r.tax_class_id),
    };
  }

  async updateTaxRate(id: number, vm: TaxRatePostVm): Promise<boolean> {
    const res = await sql`
      UPDATE tax_rate
      SET rate = ${vm.rate},
          country_id = ${vm.countryId},
          state_or_province_id = ${vm.stateOrProvinceId ?? null},
          zip_code = ${vm.zipCode ?? null},
          tax_class_id = ${vm.taxClassId}
      WHERE id = ${id}
    `;
    return res.count > 0;
  }

  async deleteTaxRate(id: number): Promise<boolean> {
    const res = await sql`
      DELETE FROM tax_rate WHERE id = ${id}
    `;
    return res.count > 0;
  }

  async getPageableTaxRates(pageNo = 0, pageSize = 10): Promise<TaxRateListVm> {
    const offset = pageNo * pageSize;
    const [countRow] = await sql`SELECT COUNT(*)::int as total FROM tax_rate`;
    const total = countRow.total || 0;

    const rows = await sql`
      SELECT tr.id, tr.rate, tr.country_id, tr.state_or_province_id, tr.zip_code, tr.tax_class_id, tc.name as tax_class_name
      FROM tax_rate tr
      LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
      ORDER BY tr.id ASC
      LIMIT ${pageSize} OFFSET ${offset}
    `;

    const totalPages = Math.ceil(total / pageSize);
    return {
      taxRates: rows.map((r: any) => ({
        id: Number(r.id),
        rate: Number(r.rate),
        countryId: Number(r.country_id),
        stateOrProvinceId: r.state_or_province_id ? Number(r.state_or_province_id) : null,
        zipCode: r.zip_code,
        taxClassId: Number(r.tax_class_id),
        taxClassName: r.tax_class_name,
      })),
      pageNo,
      pageSize,
      totalElements: total,
      totalPages,
      isLast: pageNo >= totalPages - 1,
    };
  }

  async getTaxPercent(
    taxClassId: number,
    countryId: number,
    stateOrProvinceId?: number | null,
    zipCode?: string | null
  ): Promise<number> {
    const rows = await sql`
      SELECT rate
      FROM tax_rate
      WHERE tax_class_id = ${taxClassId}
        AND country_id = ${countryId}
      ORDER BY id DESC
      LIMIT 1
    `;
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

    const rows = await sql`
      SELECT tr.id, tr.rate, tr.country_id, tr.state_or_province_id, tr.zip_code, tr.tax_class_id, tc.name as tax_class_name
      FROM tax_rate tr
      LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
      WHERE tr.country_id = ${countryId}
        AND tr.tax_class_id IN ${sql(taxClassIds)}
    `;

    return rows.map((r: any) => ({
      id: Number(r.id),
      rate: Number(r.rate),
      countryId: Number(r.country_id),
      stateOrProvinceId: r.state_or_province_id ? Number(r.state_or_province_id) : null,
      zipCode: r.zip_code,
      taxClassId: Number(r.tax_class_id),
      taxClassName: r.tax_class_name,
    }));
  }
}
