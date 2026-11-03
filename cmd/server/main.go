package main

import (
	"context"
	"io/fs"
	"log"
	"os"
	"time"

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
	r.POST("/booking", func(c *gin.Context) {
		start, err := time.Parse("2006-01-02T15:04", c.PostForm("starts_at"))
		if err != nil {
			c.String(400, "bad date")
			return
		}
		var duration int
		err = db.QueryRow(c, "SELECT duration_min FROM services WHERE id = $1", c.PostForm("service_id")).Scan(&duration)
		if err != nil {
			c.String(400, "bad service")
			return
		}
		_, err = db.Exec(c, `INSERT INTO bookings (service_id, barber_id, starts_at, ends_at, client_name, client_phone)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			c.PostForm("service_id"), c.PostForm("barber_id"), start, start.Add(time.Duration(duration)*time.Minute),
			c.PostForm("name"), c.PostForm("phone"))
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.String(200, "Спасибо! Вы записаны.")
	})
	log.Fatal(r.Run(":8080"))
}
