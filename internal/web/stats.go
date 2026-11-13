package web

import (
	"net/http"

	"barbershop/internal/store"

	"github.com/gin-gonic/gin"
)

func (s *Server) stats(c *gin.Context) {
	weeks, err := s.store.BookingsPerWeek(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	top, err := s.store.TopServices(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	s.render(c, "stats.html", gin.H{"Weeks": weeks, "WeeksMax": maxCount(weeks), "Top": top, "TopMax": maxCount(top)})
}

func maxCount(stats []store.Stat) int {
	m := 0
	for _, s := range stats {
		m = max(m, s.Count)
	}
	return m
}
