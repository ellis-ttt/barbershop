package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"barbershop/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type fakeStore struct {
	bookings []store.Booking
	blocks   []store.Block
	created  []store.Booking
	taken    bool
	hash     string
}

var services = []store.Service{{ID: 1, Name: "Стрижка", Price: 1500, DurationMin: 60}}
var barbers = []store.Barber{{ID: 1, Name: "Иван"}, {ID: 2, Name: "Пётр"}}

func (f *fakeStore) Services(context.Context) ([]store.Service, error) { return services, nil }
func (f *fakeStore) Barbers(context.Context) ([]store.Barber, error)   { return barbers, nil }

func (f *fakeStore) Service(_ context.Context, id int64) (store.Service, error) {
	if id != 1 {
		return store.Service{}, errors.New("not found")
	}
	return services[0], nil
}

func (f *fakeStore) CreateBooking(_ context.Context, b store.Booking) error {
	if f.taken {
		return store.ErrSlotTaken
	}
	f.created = append(f.created, b)
	return nil
}

func (f *fakeStore) Bookings(context.Context, time.Time, time.Time) ([]store.Booking, error) {
	return f.bookings, nil
}

func (f *fakeStore) Blocks(context.Context, time.Time, time.Time) ([]store.Block, error) {
	return f.blocks, nil
}

func (f *fakeStore) CancelBooking(context.Context, int64) error { return nil }
func (f *fakeStore) RescheduleBooking(context.Context, int64, int64, time.Time) error {
	return nil
}
func (f *fakeStore) CreateBlock(context.Context, store.Block) error { return nil }
func (f *fakeStore) DeleteBlock(context.Context, int64) error       { return nil }
func (f *fakeStore) BookingsPerWeek(context.Context) ([]store.Stat, error) {
	return []store.Stat{{Label: "02.11", Count: 3}}, nil
}
func (f *fakeStore) TopServices(context.Context) ([]store.Stat, error) {
	return []store.Stat{{Label: "Стрижка", Count: 3}}, nil
}
func (f *fakeStore) AdminHash(context.Context, string) (string, error) { return f.hash, nil }

func init() {
	gin.SetMode(gin.TestMode)
}

func do(h http.Handler, method, target string, form url.Values, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func tomorrow() time.Time {
	return time.Now().AddDate(0, 0, 1)
}

func TestIndex(t *testing.T) {
	w := do(New(&fakeStore{}, "s").Routes(), "GET", "/", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Стрижка")
	assert.Contains(t, w.Body.String(), "Пётр")
}

func TestBookingSlotsHidesBusy(t *testing.T) {
	day := tomorrow()
	start := time.Date(day.Year(), day.Month(), day.Day(), 10, 0, 0, 0, time.Local)
	st := &fakeStore{bookings: []store.Booking{{BarberID: 1, StartsAt: start, EndsAt: start.Add(time.Hour)}}}
	w := do(New(st, "s").Routes(), "GET", "/booking/slots?service_id=1&date="+day.Format(dateLayout), nil)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.NotContains(t, body, "barber_id=1&starts_at="+start.Format(dateTimeLayout))
	assert.Contains(t, body, "barber_id=2&starts_at="+start.Format(dateTimeLayout))
}

func TestCreateBooking(t *testing.T) {
	st := &fakeStore{}
	start := tomorrow().Format(dateLayout) + "T12:00"
	w := do(New(st, "s").Routes(), "POST", "/booking", url.Values{
		"service_id": {"1"}, "barber_id": {"2"}, "starts_at": {start}, "name": {"Олег"}, "phone": {"+79990000000"},
	})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Олег")
	require.Len(t, st.created, 1)
	assert.Equal(t, time.Hour, st.created[0].EndsAt.Sub(st.created[0].StartsAt))
}

func TestCreateBookingTaken(t *testing.T) {
	start := tomorrow().Format(dateLayout) + "T12:00"
	w := do(New(&fakeStore{taken: true}, "s").Routes(), "POST", "/booking", url.Values{
		"service_id": {"1"}, "barber_id": {"2"}, "starts_at": {start}, "name": {"Олег"}, "phone": {"1"},
	})
	assert.Contains(t, w.Body.String(), "уже заняли")
}

func TestAdminRequiresLogin(t *testing.T) {
	w := do(New(&fakeStore{}, "s").Routes(), "GET", "/admin", nil)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/admin/login", w.Header().Get("Location"))
}

func TestLoginAndDashboard(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	day := tomorrow()
	start := time.Date(day.Year(), day.Month(), day.Day(), 11, 0, 0, 0, time.Local)
	st := &fakeStore{
		hash:     string(hash),
		bookings: []store.Booking{{ID: 7, BarberID: 1, StartsAt: start, EndsAt: start.Add(time.Hour), ClientName: "Олег", ServiceName: "Стрижка"}},
		blocks:   []store.Block{{ID: 3, BarberID: 2, StartsAt: start, EndsAt: start.Add(time.Hour), Reason: "Обед"}},
	}
	h := New(st, "s").Routes()

	w := do(h, "POST", "/admin/login", url.Values{"login": {"admin"}, "password": {"wrong"}})
	assert.Contains(t, w.Body.String(), "Неверный")

	w = do(h, "POST", "/admin/login", url.Values{"login": {"admin"}, "password": {"secret"}})
	require.Equal(t, http.StatusSeeOther, w.Code)
	cookie := w.Result().Cookies()[0]

	w = do(h, "GET", "/admin?date="+day.Format(dateLayout), nil, cookie)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Олег")
	assert.Contains(t, w.Body.String(), "Обед")

	w = do(h, "GET", "/admin?view=week&date="+day.Format(dateLayout), nil, cookie)
	assert.Contains(t, w.Body.String(), "Олег")

	w = do(h, "GET", "/admin/stats", nil, cookie)
	assert.Contains(t, w.Body.String(), "width: 100%")
}
