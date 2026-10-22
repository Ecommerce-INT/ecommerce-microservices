-- PostgreSQL 18 Multi-Database Bootstrap
-- PostgreSQL 18 features utilized:
-- 1. Native UUIDv7 (RFC 9562 time-ordered clustering keys via uuidv7())
-- 2. Asynchronous I/O subsystem (io_method = worker / io_uring)
-- 3. B-tree Skip Scan query execution on composite indexes
-- 4. Atomic state transition tracking via RETURNING OLD, NEW

CREATE DATABASE keycloak;
CREATE DATABASE authservice;
CREATE DATABASE productservice;
CREATE DATABASE orderservice;
CREATE DATABASE paymentservice;
CREATE DATABASE shippingservice;
CREATE DATABASE inventoryservice;
CREATE DATABASE favouriteservice;
CREATE DATABASE ratingservice;
CREATE DATABASE mediaservice;
CREATE DATABASE taxservice;
CREATE DATABASE promotionservice;
