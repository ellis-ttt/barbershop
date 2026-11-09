package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) loginPage(c *gin.Context) {
	s.render(c, "login.html", nil)
}

func (s *Server) login(c *gin.Context) {
	hash, err := s.store.AdminHash(c, c.PostForm("login"))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.PostForm("password"))) != nil {
		s.render(c, "login.html", gin.H{"Error": "Неверный логин или пароль"})
		return
	}
	c.String(http.StatusOK, "ok")
}
