package config

import "os"

type Config struct {
	Addr          string
	DatabaseURL   string
	AdminLogin    string
	AdminPassword string
	JWTSecret     string
	Timezone      string
}

func Load() Config {
	return Config{
		Addr:          env("ADDR", ":8080"),
		DatabaseURL:   env("DATABASE_URL", "postgres://barber:barber@localhost:5432/barbershop"),
		AdminLogin:    env("ADMIN_LOGIN", "admin"),
		AdminPassword: env("ADMIN_PASSWORD", "admin"),
		JWTSecret:     env("JWT_SECRET", "dev-secret"),
		Timezone:      env("TZ", "Europe/Moscow"),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
