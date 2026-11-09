package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenCookie = "token"
	tokenTTL    = 24 * time.Hour
)

func (s *Server) loginPage(c *gin.Context) {
	s.render(c, "login.html", nil)
}

func (s *Server) login(c *gin.Context) {
	login := c.PostForm("login")
	hash, err := s.store.AdminHash(c, login)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.PostForm("password"))) != nil {
		s.render(c, "login.html", gin.H{"Error": "Неверный логин или пароль"})
		return
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	}).SignedString(s.secret)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(tokenCookie, token, int(tokenTTL.Seconds()), "/admin", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin")
}

func (s *Server) logout(c *gin.Context) {
	c.SetCookie(tokenCookie, "", -1, "/admin", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func (s *Server) requireAdmin(c *gin.Context) {
	cookie, err := c.Cookie(tokenCookie)
	if err == nil {
		_, err = jwt.Parse(cookie, func(*jwt.Token) (any, error) { return s.secret, nil },
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	}
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		c.Abort()
	}
}
