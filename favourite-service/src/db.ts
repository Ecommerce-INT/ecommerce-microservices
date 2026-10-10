import postgres from "postgres";

function getDatabaseUrl(): string {
  if (process.env.DATABASE_URL) {
    return process.env.DATABASE_URL;
  }

  const user = process.env.POSTGRES_USER || "postgres";
  const pass = process.env.POSTGRES_PASSWORD || "postgres";
  let host = process.env.POSTGRES_HOST || "localhost";
  let port = process.env.POSTGRES_PORT || "5432";
  let db = process.env.POSTGRES_DB || "favouriteservice";

  // Support SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/favouriteservice
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
      CREATE TABLE IF NOT EXISTS favourites (
        user_id INT NOT NULL,
        product_id INT NOT NULL,
        like_date TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (user_id, product_id, like_date)
      );
    `;
    console.log("[Favourite Service] Database table 'favourites' ready.");
  } catch (err) {
    console.warn("[Favourite Service] Notice during table verification:", err);
  }
}
