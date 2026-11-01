package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Бритва — скоро открытие")
	})
	log.Fatal(r.Run(":8080"))
}
