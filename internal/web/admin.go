package web

import "github.com/gin-gonic/gin"

func (s *Server) dashboard(c *gin.Context) {
	s.render(c, "dashboard.html", nil)
}
