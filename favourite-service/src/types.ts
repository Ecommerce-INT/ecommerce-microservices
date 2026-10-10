import { z } from "zod";

export const FavouriteIdSchema = z.object({
  userId: z.number().int(),
  productId: z.number().int(),
  likeDate: z.string(),
});

export type FavouriteId = z.infer<typeof FavouriteIdSchema>;

export const FavouriteDtoSchema = z.object({
  userId: z.number().int(),
  productId: z.number().int(),
  likeDate: z.string(),
  user: z.any().optional(),
  product: z.any().optional(),
});

export type FavouriteDto = z.infer<typeof FavouriteDtoSchema>;

export interface DtoCollectionResponse<T> {
  collection: T[];
}
