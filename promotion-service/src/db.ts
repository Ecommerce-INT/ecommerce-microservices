import postgres from "postgres";

function getDatabaseUrl(): string {
  if (process.env.DATABASE_URL) {
    return process.env.DATABASE_URL;
  }

  const user = process.env.POSTGRES_USER || "postgres";
  const pass = process.env.POSTGRES_PASSWORD || "postgres";
  let host = process.env.POSTGRES_HOST || "localhost";
  let port = process.env.POSTGRES_PORT || "5432";
  let db = process.env.POSTGRES_DB || "promotionservice";

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
      CREATE TABLE IF NOT EXISTS promotion (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        slug VARCHAR(255) NOT NULL,
        description VARCHAR(255),
        coupon_code VARCHAR(255),
        discount_percentage BIGINT NOT NULL DEFAULT 0,
        discount_amount BIGINT NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        start_date TIMESTAMP WITH TIME ZONE,
        end_date TIMESTAMP WITH TIME ZONE,
        discount_type VARCHAR(25) NOT NULL DEFAULT 'PERCENTAGE',
        usage_limit INT DEFAULT 100,
        usage_count INT DEFAULT 0 NOT NULL,
        usage_type VARCHAR(25) DEFAULT 'UNLIMITED' NOT NULL,
        apply_to VARCHAR(25) DEFAULT 'ALL' NOT NULL,
        minimum_order_purchase_amount BIGINT DEFAULT 0,
        created_by VARCHAR(255),
        created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        last_modified_by VARCHAR(255),
        last_modified_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
      );
    `;
    console.log("[Promotion Service] Database table 'promotion' ready.");
  } catch (err) {
    console.warn("[Promotion Service] Notice during database init:", err);
  }
}
