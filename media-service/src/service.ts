import { sql } from "./db";
import { uploadToStorage, deleteFromStorage, getFromStorage } from "./storage";

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

const publicUrl = process.env.PUBLIC_URL || "http://localhost:8083";

function buildMediaUrl(id: number, fileName: string): string {
  return `${publicUrl}/medias/${id}/file/${encodeURIComponent(fileName)}`;
}

function buildObjectKey(fileName: string): string {
  const safeName = (fileName || "file").replace(/[^a-zA-Z0-9._-]/g, "_");
  return `${crypto.randomUUID()}-${safeName}`;
}

export class MediaService {
  async saveMedia(file: File, caption = "", fileNameOverride = ""): Promise<NoFileMediaVm> {
    const fileName = fileNameOverride.trim() || file.name || "upload.bin";
    const mediaType = file.type || "application/octet-stream";
    const objectKey = buildObjectKey(fileName);

    const arrayBuffer = await file.arrayBuffer();
    await uploadToStorage(objectKey, new Uint8Array(arrayBuffer), mediaType);

    const rows = await sql`
      INSERT INTO media (caption, file_name, file_path, media_type)
      VALUES (${caption}, ${fileName}, ${objectKey}, ${mediaType})
      RETURNING id, caption, file_name, media_type
    `;

    const r = rows[0];
    return {
      id: Number(r.id),
      caption: r.caption || "",
      fileName: r.file_name || "",
      mediaType: r.media_type || "",
    };
  }

  async removeMedia(id: number): Promise<void> {
    const rows = await sql`SELECT file_path FROM media WHERE id = ${id}`;
    if (rows.length > 0 && rows[0].file_path) {
      await deleteFromStorage(rows[0].file_path);
    }
    await sql`DELETE FROM media WHERE id = ${id}`;
  }

  async getMediaById(id: number): Promise<MediaVm | null> {
    const rows = await sql`SELECT id, caption, file_name, media_type FROM media WHERE id = ${id}`;
    if (rows.length === 0) {
      return null;
    }
    const r = rows[0];
    const mediaId = Number(r.id);
    const fileName = r.file_name || "";
    return {
      id: mediaId,
      caption: r.caption || "",
      fileName,
      mediaType: r.media_type || "",
      url: buildMediaUrl(mediaId, fileName),
    };
  }

  async getMediaByIds(ids: number[]): Promise<MediaVm[]> {
    if (ids.length === 0) return [];
    const rows = await sql`SELECT id, caption, file_name, media_type FROM media WHERE id = ANY(${ids})`;
    return rows.map((r: any) => {
      const mediaId = Number(r.id);
      const fileName = r.file_name || "";
      return {
        id: mediaId,
        caption: r.caption || "",
        fileName,
        mediaType: r.media_type || "",
        url: buildMediaUrl(mediaId, fileName),
      };
    });
  }

  async getFile(id: number, fileName: string): Promise<{ stream: any; mediaType: string } | null> {
    const rows = await sql`SELECT file_path, file_name, media_type FROM media WHERE id = ${id}`;
    if (rows.length === 0) return null;
    const r = rows[0];
    if (r.file_name && r.file_name.toLowerCase() !== fileName.toLowerCase()) {
      return null;
    }
    const stream = await getFromStorage(r.file_path);
    return {
      stream,
      mediaType: r.media_type || "application/octet-stream",
    };
  }
}
