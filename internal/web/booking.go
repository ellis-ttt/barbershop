package web

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"

	"barbershop/internal/schedule"
	"barbershop/internal/store"

	"github.com/gin-gonic/gin"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02T15:04"
)

type cell struct {
	BarberID int64
	Time     time.Time
	Free     bool
}

type row struct {
	Time  time.Time
	Cells []cell
}

func (s *Server) service(c *gin.Context) (store.Service, bool) {
	id, _ := strconv.ParseInt(c.Query("service_id"), 10, 64)
	service, err := s.store.Service(c, id)
	if err != nil {
		c.String(http.StatusBadRequest, "Услуга не найдена")
		return service, false
	}
	return service, true
}

func (s *Server) bookingDate(c *gin.Context) {
	service, ok := s.service(c)
	if !ok {
		return
	}
	s.partial(c, "booking-date", gin.H{"Service": service})
}

func (s *Server) bookingSlots(c *gin.Context) {
	service, ok := s.service(c)
	if !ok {
		return
	}
	day, err := time.ParseInLocation(dateLayout, c.Query("date"), time.Local)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверная дата")
		return
	}
	barbers, err := s.store.Barbers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	bookings, err := s.store.Bookings(c, day, day.AddDate(0, 0, 1))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	blocks, err := s.store.Blocks(c, day, day.AddDate(0, 0, 1))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	busy := map[int64][]schedule.Interval{}
	for _, b := range bookings {
		busy[b.BarberID] = append(busy[b.BarberID], schedule.Interval{Start: b.StartsAt, End: b.EndsAt})
	}
	for _, b := range blocks {
		busy[b.BarberID] = append(busy[b.BarberID], schedule.Interval{Start: b.StartsAt, End: b.EndsAt})
	}
	free := map[int64][]time.Time{}
	for _, b := range barbers {
		free[b.ID] = schedule.Free(day, service.Duration(), busy[b.ID], time.Now())
	}
	var rows []row
	for _, t := range schedule.Slots(day) {
		r := row{Time: t}
		for _, b := range barbers {
			r.Cells = append(r.Cells, cell{BarberID: b.ID, Time: t, Free: slices.Contains(free[b.ID], t)})
		}
		rows = append(rows, r)
	}
	s.partial(c, "booking-slots", gin.H{"Service": service, "Barbers": barbers, "Rows": rows})
}

func (s *Server) bookingDetails(c *gin.Context) {
	service, ok := s.service(c)
	if !ok {
		return
	}
	start, err := time.ParseInLocation(dateTimeLayout, c.Query("starts_at"), time.Local)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверное время")
		return
	}
	s.partial(c, "booking-details", gin.H{"Service": service, "BarberID": c.Query("barber_id"), "Start": start})
}

func (s *Server) createBooking(c *gin.Context) {
	start, err := time.ParseInLocation(dateTimeLayout, c.PostForm("starts_at"), time.Local)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверное время")
		return
	}
	serviceID, _ := strconv.ParseInt(c.PostForm("service_id"), 10, 64)
	barberID, _ := strconv.ParseInt(c.PostForm("barber_id"), 10, 64)
	service, err := s.store.Service(c, serviceID)
	if err != nil {
		c.String(http.StatusBadRequest, "Услуга не найдена")
		return
	}
	b := store.Booking{
		ServiceID:   serviceID,
		BarberID:    barberID,
		StartsAt:    start,
		EndsAt:      start.Add(service.Duration()),
		ClientName:  c.PostForm("name"),
		ClientPhone: c.PostForm("phone"),
	}
	err = s.store.CreateBooking(c, b)
	if errors.Is(err, store.ErrSlotTaken) {
		s.partial(c, "booking-taken", gin.H{"Service": service})
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	s.partial(c, "booking-done", gin.H{"Service": service, "Booking": b})
}
