import { z } from "zod";

export const TaxClassPostVmSchema = z.object({
  name: z.string().min(1),
});

export type TaxClassPostVm = z.infer<typeof TaxClassPostVmSchema>;

export interface TaxClassVm {
  id: number;
  name: string;
}

export interface TaxClassListVm {
  taxClasses: TaxClassVm[];
  pageNo: number;
  pageSize: number;
  totalElements: number;
  totalPages: number;
  isLast: boolean;
}

export const TaxRatePostVmSchema = z.object({
  rate: z.number(),
  countryId: z.number().int(),
  stateOrProvinceId: z.number().int().optional().nullable(),
  zipCode: z.string().optional().nullable(),
  taxClassId: z.number().int(),
});

export type TaxRatePostVm = z.infer<typeof TaxRatePostVmSchema>;

export interface TaxRateVm {
  id: number;
  rate: number;
  countryId: number;
  stateOrProvinceId?: number | null;
  zipCode?: string | null;
  taxClassId: number;
  taxClassName?: string | null;
}

export interface TaxRateListVm {
  taxRates: TaxRateVm[];
  pageNo: number;
  pageSize: number;
  totalElements: number;
  totalPages: number;
  isLast: boolean;
}
