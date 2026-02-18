-- ============================================================
-- E-Commerce Platform — PostgreSQL Init Script
-- Her mikroservis için ayrı veritabanı oluşturur
-- Dosya yolu: init-scripts/postgres/01-init-databases.sql
-- ============================================================

-- user-service veritabanı
CREATE DATABASE users_db;
GRANT ALL PRIVILEGES ON DATABASE users_db TO ecommerce;

-- product-service veritabanı
CREATE DATABASE products_db;
GRANT ALL PRIVILEGES ON DATABASE products_db TO ecommerce;

-- order-service veritabanı
CREATE DATABASE orders_db;
GRANT ALL PRIVILEGES ON DATABASE orders_db TO ecommerce;

-- payment-service veritabanı
CREATE DATABASE payments_db;
GRANT ALL PRIVILEGES ON DATABASE payments_db TO ecommerce;

-- notification-service veritabanı
CREATE DATABASE notifications_db;
GRANT ALL PRIVILEGES ON DATABASE notifications_db TO ecommerce;

-- Keycloak veritabanı
CREATE DATABASE keycloak_db;
GRANT ALL PRIVILEGES ON DATABASE keycloak_db TO ecommerce;
