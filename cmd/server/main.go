package main

import (
	"context"
	"log"
	"time"
	_ "time/tzdata"

	"barbershop/internal/config"
	"barbershop/internal/store"
	"barbershop/internal/web"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc
	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	if err := st.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	if err := st.SaveAdmin(ctx, cfg.AdminLogin, string(hash)); err != nil {
		log.Fatal(err)
	}
	log.Fatal(web.New(st, cfg.JWTSecret).Routes().Run(cfg.Addr))
}
