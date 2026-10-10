import { z } from "zod";

export const PromotionPostVmSchema = z.object({
  name: z.string().min(1),
  slug: z.string().optional(),
  description: z.string().optional().default(""),
  couponCode: z.string().optional().default(""),
  usageLimit: z.number().int().optional().default(100),
  discountType: z.string().optional().default("PERCENTAGE"),
  applyTo: z.string().optional().default("ALL"),
  usageType: z.string().optional().default("UNLIMITED"),
  discountPercentage: z.number().int().optional().default(0),
  discountAmount: z.number().int().optional().default(0),
  isActive: z.boolean().optional().default(true),
  startDate: z.string().optional(),
  endDate: z.string().optional(),
});

export type PromotionPostVm = z.infer<typeof PromotionPostVmSchema>;

export const PromotionPutVmSchema = PromotionPostVmSchema.extend({
  id: z.number().int(),
});

export type PromotionPutVm = z.infer<typeof PromotionPutVmSchema>;

export interface PromotionDetailVm {
  id: number;
  name: string;
  slug: string;
  description: string;
  couponCode: string;
  usageLimit: number;
  usageCount: number;
  discountType: string;
  applyTo: string;
  usageType: string;
  discountPercentage: number;
  discountAmount: number;
  isActive: boolean;
  startDate: string | null;
  endDate: string | null;
  brands: any[];
  categories: any[];
  products: any[];
}

export interface PromotionListVm {
  promotionDetailVmList: PromotionDetailVm[];
  totalElements: number;
  totalPages: number;
}
