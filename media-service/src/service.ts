import { eq, inArray } from 'drizzle-orm';
import { db } from './db';
import { media } from './schema';
import { uploadToStorage, deleteFromStorage, getFromStorage } from './storage';

export interface MediaVm {
  id: number;
  caption: string;
  fileName: string;
  mediaType: string;
  url: string;
}

export interface NoFileMediaVm {
  id: number;
  caption: string;
  fileName: string;
  mediaType: string;
}

const publicUrl = process.env.PUBLIC_URL || 'http://localhost:8083';

function buildMediaUrl(id: number, fileName: string): string {
  return `${publicUrl}/medias/${id}/file/${encodeURIComponent(fileName)}`;
}

function buildObjectKey(fileName: string): string {
  const safeName = (fileName || 'file').replace(/[^a-zA-Z0-9._-]/g, '_');
  return `${crypto.randomUUID()}-${safeName}`;
}

export class MediaService {
  async saveMedia(file: File, caption = '', fileNameOverride = ''): Promise<NoFileMediaVm> {
    const fileName = fileNameOverride.trim() || file.name || 'upload.bin';
    const mediaType = file.type || 'application/octet-stream';
    const objectKey = buildObjectKey(fileName);

    const arrayBuffer = await file.arrayBuffer();
    await uploadToStorage(objectKey, new Uint8Array(arrayBuffer), mediaType);

    const [row] = await db
      .insert(media)
      .values({
        caption,
        fileName,
        filePath: objectKey,
        mediaType
      })
      .returning();

    return {
      id: row.id,
      caption: row.caption || '',
      fileName: row.fileName || '',
      mediaType: row.mediaType || ''
    };
  }

  async removeMedia(id: number): Promise<void> {
    const rows = await db
      .select({ filePath: media.filePath })
      .from(media)
      .where(eq(media.id, id));
    if (rows.length > 0 && rows[0].filePath) {
      await deleteFromStorage(rows[0].filePath);
    }
    await db.delete(media).where(eq(media.id, id));
  }

  async getMediaById(id: number): Promise<MediaVm | null> {
    const rows = await db.select().from(media).where(eq(media.id, id));
    if (rows.length === 0) {
      return null;
    }
    const r = rows[0];
    const fileName = r.fileName || '';
    return {
      id: r.id,
      caption: r.caption || '',
      fileName,
      mediaType: r.mediaType || '',
      url: buildMediaUrl(r.id, fileName)
    };
  }

  async getMediaByIds(ids: number[]): Promise<MediaVm[]> {
    if (ids.length === 0) return [];
    const rows = await db.select().from(media).where(inArray(media.id, ids));
    return rows.map((r) => {
      const fileName = r.fileName || '';
      return {
        id: r.id,
        caption: r.caption || '',
        fileName,
        mediaType: r.mediaType || '',
        url: buildMediaUrl(r.id, fileName)
      };
    });
  }

  async getFile(
    id: number,
    fileName: string
  ): Promise<{ stream: ReadableStream<Uint8Array>; mediaType: string } | null> {
    const rows = await db.select().from(media).where(eq(media.id, id));
    if (rows.length === 0) return null;
    const r = rows[0];
    if (r.fileName && r.fileName.toLowerCase() !== fileName.toLowerCase()) {
      return null;
    }
    if (!r.filePath) return null;
    const stream = await getFromStorage(r.filePath);
    if (!stream) return null;
    return {
      stream,
      mediaType: r.mediaType || 'application/octet-stream'
    };
  }
}
