import { S3Client } from "bun";

const endpoint = process.env.STORAGE_ENDPOINT || "http://localhost:9000";
const region = process.env.STORAGE_REGION || "us-east-1";
const accessKeyId = process.env.STORAGE_ACCESS_KEY || "admin";
const secretAccessKey = process.env.STORAGE_SECRET_KEY || "admin";
export const bucket = process.env.STORAGE_BUCKET || "ecommerce-media";

export const s3 = new S3Client({
  endpoint,
  region,
  accessKeyId,
  secretAccessKey,
  bucket,
  virtualHostedStyle: false,
});

export async function initStorage() {
  try {
    await s3.list({ maxKeys: 1 });
    console.log(`[Media Service] Bucket '${bucket}' ready on ${endpoint}`);
  } catch (err: any) {
    console.warn(
      `[Media Service] Bucket '${bucket}' is not reachable on ${endpoint}:`,
      err?.message ?? err
    );
    console.warn(
      `[Media Service] Create it once (docker compose up storage-init) before uploading.`
    );
  }
}

export async function uploadToStorage(key: string, body: Uint8Array | Buffer, contentType: string) {
  await s3.write(key, body, { type: contentType });
}

export async function deleteFromStorage(key: string) {
  try {
    await s3.unlink(key);
  } catch (err: any) {
    console.warn(`[Media Service] Failed to delete object ${key}:`, err.message);
  }
}

export async function getFromStorage(key: string): Promise<ReadableStream<Uint8Array> | null> {
  const file = s3.file(key);
  if (!(await file.exists())) {
    return null;
  }
  return file.stream();
}
