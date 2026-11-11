package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"barbershop/internal/schedule"
	"barbershop/internal/store"

	"github.com/gin-gonic/gin"
)

type adminCell struct {
	Booking *store.Booking
	Start   bool
}

type adminRow struct {
	Time  time.Time
	Cells []adminCell
}

type dayBookings struct {
	Day      time.Time
	Bookings []store.Booking
}

func (s *Server) dashboard(c *gin.Context) {
	day, err := time.Parse(dateLayout, c.Query("date"))
	if err != nil {
		y, m, d := time.Now().Date()
		day = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}
	week := c.Query("view") == "week"
	to := day.AddDate(0, 0, 1)
	if week {
		to = day.AddDate(0, 0, 7)
	}
	barbers, err := s.store.Barbers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	bookings, err := s.store.Bookings(c, day, to)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	data := gin.H{"Day": day, "Week": week, "Barbers": barbers, "Prev": day.AddDate(0, 0, -1), "Next": day.AddDate(0, 0, 1)}
	if week {
		data["Prev"], data["Next"] = day.AddDate(0, 0, -7), to
		data["Days"] = groupByDay(day, bookings)
	} else {
		data["Rows"] = dayGrid(day, barbers, bookings)
	}
	s.render(c, "dashboard.html", data)
}

func dayGrid(day time.Time, barbers []store.Barber, bookings []store.Booking) []adminRow {
	var rows []adminRow
	for _, t := range schedule.Slots(day) {
		r := adminRow{Time: t}
		for _, br := range barbers {
			var cell adminCell
			for i, b := range bookings {
				if b.BarberID == br.ID && !t.Before(b.StartsAt) && t.Before(b.EndsAt) {
					cell = adminCell{Booking: &bookings[i], Start: t.Equal(b.StartsAt)}
				}
			}
			r.Cells = append(r.Cells, cell)
		}
		rows = append(rows, r)
	}
	return rows
}

func groupByDay(from time.Time, bookings []store.Booking) []dayBookings {
	days := make([]dayBookings, 7)
	for i := range days {
		days[i].Day = from.AddDate(0, 0, i)
	}
	for _, b := range bookings {
		i := int(b.StartsAt.Sub(from).Hours() / 24)
		days[i].Bookings = append(days[i].Bookings, b)
	}
	return days
}

func (s *Server) cancelBooking(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := s.store.CancelBooking(c, id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}

func (s *Server) rescheduleBooking(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	barberID, _ := strconv.ParseInt(c.PostForm("barber_id"), 10, 64)
	start, err := time.Parse(dateTimeLayout, c.PostForm("starts_at"))
	if err != nil {
		c.String(http.StatusBadRequest, "Неверное время")
		return
	}
	err = s.store.RescheduleBooking(c, id, barberID, start)
	if errors.Is(err, store.ErrSlotTaken) {
		c.String(http.StatusConflict, "Это время занято")
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin?date="+start.Format(dateLayout))
}
