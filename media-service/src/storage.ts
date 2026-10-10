import {
  S3Client,
  PutObjectCommand,
  GetObjectCommand,
  DeleteObjectCommand,
  CreateBucketCommand,
} from "@aws-sdk/client-s3";

const endpoint = process.env.STORAGE_ENDPOINT || "http://localhost:9000";
const region = process.env.STORAGE_REGION || "us-east-1";
const accessKeyId = process.env.STORAGE_ACCESS_KEY || "admin";
const secretAccessKey = process.env.STORAGE_SECRET_KEY || "admin";
export const bucket = process.env.STORAGE_BUCKET || "ecommerce-media";

export const s3 = new S3Client({
  endpoint,
  region,
  credentials: {
    accessKeyId,
    secretAccessKey,
  },
  forcePathStyle: true,
});

export async function initStorage() {
  try {
    await s3.send(new CreateBucketCommand({ Bucket: bucket }));
    console.log(`[Media Service] Bucket '${bucket}' ensured on ${endpoint}`);
  } catch (err: any) {
    if (err.name === "BucketAlreadyOwnedByYou" || err.name === "BucketAlreadyExists") {
      // already exists, all good
    } else {
      console.warn(`[Media Service] Notice during S3 init (${endpoint}):`, err.message);
    }
  }
}

export async function uploadToStorage(key: string, body: Uint8Array | Buffer, contentType: string) {
  await s3.send(
    new PutObjectCommand({
      Bucket: bucket,
      Key: key,
      Body: body,
      ContentType: contentType,
    })
  );
}

export async function deleteFromStorage(key: string) {
  try {
    await s3.send(
      new DeleteObjectCommand({
        Bucket: bucket,
        Key: key,
      })
    );
  } catch (err: any) {
    console.warn(`[Media Service] Failed to delete object ${key}:`, err.message);
  }
}

export async function getFromStorage(key: string) {
  const res = await s3.send(
    new GetObjectCommand({
      Bucket: bucket,
      Key: key,
    })
  );
  return res.Body;
}
