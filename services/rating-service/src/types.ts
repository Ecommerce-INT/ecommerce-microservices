import { z } from "zod";

export const RatingPostVmSchema = z.object({
  content: z.string().optional().default(""),
  star: z.number().int().min(1).max(5),
  productId: z.number().int(),
  productName: z.string().optional().default(""),
});

export type RatingPostVm = z.infer<typeof RatingPostVmSchema>;

export interface RatingVm {
  id: number;
  content: string;
  star: number;
  productId: number;
  productName: string;
  createdBy: string;
  lastName: string;
  firstName: string;
  createdOn: string;
}

export interface RatingListVm {
  ratingList: RatingVm[];
  totalElements: number;
  totalPages: number;
}

export interface ResponeStatusVm {
  title: string;
  message: string;
  statusCode: string;
}
