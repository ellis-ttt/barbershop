package web

import (
	"net/http"
	"strconv"
	"time"

	"barbershop/internal/store"

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

func (s *Server) createBooking(c *gin.Context) {
	start, err := time.Parse("2006-01-02T15:04", c.PostForm("starts_at"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad date")
		return
	}
	serviceID, _ := strconv.ParseInt(c.PostForm("service_id"), 10, 64)
	barberID, _ := strconv.ParseInt(c.PostForm("barber_id"), 10, 64)
	service, err := s.store.Service(c, serviceID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad service")
		return
	}
	err = s.store.CreateBooking(c, store.Booking{
		ServiceID:   serviceID,
		BarberID:    barberID,
		StartsAt:    start,
		EndsAt:      start.Add(service.Duration()),
		ClientName:  c.PostForm("name"),
		ClientPhone: c.PostForm("phone"),
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "Спасибо! Вы записаны.")
}
