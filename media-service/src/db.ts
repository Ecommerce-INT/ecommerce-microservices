import postgres from "postgres";

function getDatabaseUrl(): string {
  if (process.env.DATABASE_URL) {
    return process.env.DATABASE_URL;
  }

  const user = process.env.POSTGRES_USER || "postgres";
  const pass = process.env.POSTGRES_PASSWORD || "postgres";
  let host = process.env.POSTGRES_HOST || "localhost";
  let port = process.env.POSTGRES_PORT || "5432";
  let db = process.env.POSTGRES_DB || "mediaservice";

  if (process.env.SPRING_DATASOURCE_URL) {
    const raw = process.env.SPRING_DATASOURCE_URL.replace("jdbc:postgresql://", "");
    const parts = raw.split("/");
    if (parts.length >= 2) {
      const [hostPort, dbName] = parts;
      db = dbName;
      if (hostPort.includes(":")) {
        const [h, p] = hostPort.split(":");
        host = h;
        port = p;
      } else {
        host = hostPort;
      }
    }
  }

  return `postgres://${user}:${pass}@${host}:${port}/${db}?sslmode=disable`;
}

const connectionString = getDatabaseUrl();

export const sql = postgres(connectionString, {
  max: 20,
  idle_timeout: 30,
  connect_timeout: 5,
  onnotice: () => {},
});

export async function initDb() {
  try {
    await sql`
      CREATE TABLE IF NOT EXISTS media (
        id BIGSERIAL PRIMARY KEY,
        caption VARCHAR(255),
        file_name VARCHAR(255),
        file_path TEXT,
        media_type VARCHAR(128)
      );
    `;
    console.log("[Media Service] Database table 'media' ready.");
  } catch (err) {
    console.warn("[Media Service] Notice during database init:", err);
  }
}
