package main

import (
	"context"
	"log"
	"os"

	"barbershop/internal/store"
	"barbershop/internal/web"
)

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://barber:barber@localhost:5432/barbershop"
	}
	st, err := store.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	if err := st.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	log.Fatal(web.New(st).Routes().Run(":8080"))
}
