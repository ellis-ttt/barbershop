package main

import (
	"context"
	"io/fs"
	"log"
	"os"

	"barbershop/migrations"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	ID          int64
	Name        string
	Price       int
	DurationMin int
}

type Barber struct {
	ID   int64
	Name string
}

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://barber:barber@localhost:5432/barbershop"
	}
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, _ := fs.Glob(migrations.FS, "*.sql")
	for _, f := range files {
		sql, _ := migrations.FS.ReadFile(f)
		if _, err := db.Exec(ctx, string(sql)); err != nil {
			log.Fatal(f, err)
		}
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		rows, _ := db.Query(c, "SELECT id, name, price, duration_min FROM services ORDER BY id")
		services, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Service])
		if err != nil {
			c.String(500, err.Error())
			return
		}
		rows, _ = db.Query(c, "SELECT id, name FROM barbers ORDER BY id")
		barbers, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Barber])
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.HTML(200, "index.html", gin.H{"Services": services, "Barbers": barbers})
	})
	log.Fatal(r.Run(":8080"))
}
