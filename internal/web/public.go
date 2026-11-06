package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) index(c *gin.Context) {
	services, err := s.store.Services(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	barbers, err := s.store.Barbers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	s.render(c, "index.html", gin.H{"Services": services, "Barbers": barbers})
}
