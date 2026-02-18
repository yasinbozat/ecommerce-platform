# E-Commerce Platform — Infrastructure

## Klasör Yapısı

```
infra/
├── docker-compose.yml
├── init-scripts/
│   └── postgres/
│       └── 01-init-databases.sql     # DB'leri otomatik oluşturur
├── nginx/
│   ├── nginx.conf                    # Ana Nginx config
│   └── conf.d/
│       └── ecommerce.conf            # Routing + auth_request
├── prometheus/
│   └── prometheus.yml                # Scrape config
├── grafana/
│   ├── provisioning/                 # Datasource otomatik bağlama
│   └── dashboards/                   # Dashboard JSON'ları
├── pgadmin/
│   └── servers.json                  # pgAdmin otomatik bağlantı
└── keycloak/
    └── realm-export.json             # Realm otomatik import
```

## Başlatma

```bash
# Tüm altyapıyı ayağa kaldır
docker-compose up -d

# Logları takip et
docker-compose logs -f

# Belirli bir servisin logları
docker-compose logs -f kafka

# Durdur
docker-compose down

# Durdur ve volume'ları sil (temiz başlangıç)
docker-compose down -v
```

## Servis Adresleri

| Servis           | URL                          | Kimlik Bilgileri              |
|------------------|------------------------------|-------------------------------|
| Nginx            | http://localhost:80          | —                             |
| Keycloak Admin   | http://localhost:8180        | admin / admin123              |
| Kafdrop          | http://localhost:9000        | —                             |
| Kibana           | http://localhost:5601        | —                             |
| Grafana          | http://localhost:3000        | admin / grafana123            |
| Jaeger UI        | http://localhost:16686       | —                             |
| pgAdmin          | http://localhost:5050        | admin@ecommerce.com / pgadmin123 |
| PostgreSQL       | localhost:5432               | ecommerce / ecommerce123      |
| Redis            | localhost:6379               | redis123                      |
| Kafka            | localhost:9092               | —                             |
| Elasticsearch    | http://localhost:9200        | —                             |

## Veritabanları

PostgreSQL içinde otomatik oluşturulan veritabanları:

| Veritabanı       | Servis                |
|------------------|-----------------------|
| users_db         | user-service          |
| products_db      | product-service       |
| orders_db        | order-service         |
| payments_db      | payment-service       |
| notifications_db | notification-service  |
| keycloak_db      | Keycloak              |

## Kafka Topikleri

Otomatik oluşturulan topikler:

- product.created / product.updated / product.deleted
- order.created / order.confirmed / order.cancelled
- stock.reserve / stock.reserved / stock.reserve.failed / stock.release
- payment.process / payment.completed / payment.failed
- notification.send

## Keycloak Kurulum Sonrası

1. http://localhost:8180/admin adresine git (admin / admin123)
2. `ecommerce` realm otomatik import edilmiş olmalı
3. Yoksa: Realm → Create → import `keycloak/realm-export.json`
4. `ecommerce-app` client'ı frontend için kullan (public)
5. `ecommerce-service` client'ı backend servisler için kullan (secret gerekir)

## Nginx + Keycloak Auth Akışı

```
Client
  → Nginx :80
  → auth_request → user-service:8081/internal/auth/validate
      → Keycloak:8180/realms/ecommerce/protocol/openid-connect/token/introspect
  → Downstream servis (X-User-ID, X-User-Role header ile)
```

## Servislerini Eklerken

Servislerini yazdıktan sonra docker-compose.yml'e ekle:

```yaml
user-service:
  build: ../services/user-service
  container_name: ecommerce-user-service
  ports:
    - "8081:8081"
  environment:
    DB_HOST: postgres
    DB_PORT: 5432
    DB_NAME: users_db
    DB_USER: ecommerce
    DB_PASSWORD: ecommerce123
    REDIS_ADDR: redis:6379
    REDIS_PASSWORD: redis123
    KEYCLOAK_URL: http://keycloak:8180
    KEYCLOAK_REALM: ecommerce
    JAEGER_ENDPOINT: http://jaeger:4318/v1/traces
  depends_on:
    postgres:
      condition: service_healthy
    redis:
      condition: service_healthy
    keycloak:
      condition: service_healthy
  networks:
    - ecommerce-net
```
