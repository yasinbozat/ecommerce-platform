package config

import "os"

type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	Keycloak KeycloakConfig
	Jaeger   JaegerConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       string
}

type KeycloakConfig struct {
	Url          string
	Realm        string
	ClientID     string
	ClientSecret string
}

type JaegerConfig struct {
	Endpoint string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8081"),
			Env:  getEnv("APP_ENV", "development"),
		}, DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "users_db"),
			User:     getEnv("DB_USER", "ecommerce"),
			Password: getEnv("DB_PASSWORD", "ecommerce123"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}, Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "redis123"),
			DB:       getEnv("REDIS_DB", "0")},
		Keycloak: KeycloakConfig{
			Url:          getEnv("KEYCLOAK_URL", "http://localhost:8180"),
			Realm:        getEnv("KEYCLOAK_REALM", "ecommerce"),
			ClientID:     getEnv("KEYCLOAK_CLIENT_ID", "ecommerce-service"),
			ClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", "service-secret"),
		},
		Jaeger: JaegerConfig{Endpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:4318/v1/traces")},
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
