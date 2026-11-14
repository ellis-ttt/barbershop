package web

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"time"
	"unicode/utf8"

	"barbershop/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//go:embed templates static
var files embed.FS

type Store interface {
	Services(ctx context.Context) ([]store.Service, error)
	Service(ctx context.Context, id int64) (store.Service, error)
	Barbers(ctx context.Context) ([]store.Barber, error)
	CreateBooking(ctx context.Context, b store.Booking) error
	Bookings(ctx context.Context, from, to time.Time) ([]store.Booking, error)
	CancelBooking(ctx context.Context, id int64) error
	RescheduleBooking(ctx context.Context, id, barberID int64, start time.Time) error
	Blocks(ctx context.Context, from, to time.Time) ([]store.Block, error)
	CreateBlock(ctx context.Context, b store.Block) error
	DeleteBlock(ctx context.Context, id int64) error
	BookingsPerWeek(ctx context.Context) ([]store.Stat, error)
	TopServices(ctx context.Context) ([]store.Stat, error)
	AdminHash(ctx context.Context, login string) (string, error)
}

type Server struct {
	store    Store
	partials *template.Template
	pages    map[string]*template.Template
	secret   []byte
	gallery  []string
}

func New(st Store, secret string) *Server {
	partials := template.Must(template.New("").Funcs(funcs).ParseFS(files, "templates/layout.html", "templates/partials/*.html"))
	return &Server{store: st, partials: partials, pages: loadPages(partials), secret: []byte(secret), gallery: gallery()}
}

func gallery() []string {
	entries, _ := fs.ReadDir(files, "static/gallery")
	var images []string
	for _, e := range entries {
		images = append(images, "/static/gallery/"+e.Name())
	}
	return images
}

func (s *Server) Routes() *gin.Engine {
	r := gin.Default()
	static, _ := fs.Sub(files, "static")
	r.StaticFS("/static", http.FS(static))
	r.GET("/", s.index)
	r.GET("/booking/date", s.bookingDate)
	r.GET("/booking/slots", s.bookingSlots)
	r.GET("/booking/details", s.bookingDetails)
	r.POST("/booking", s.createBooking)
	r.GET("/admin/login", s.loginPage)
	r.POST("/admin/login", s.login)
	r.POST("/admin/logout", s.logout)

	admin := r.Group("/admin", s.requireAdmin)
	admin.GET("", s.dashboard)
	admin.POST("/bookings/:id/cancel", s.cancelBooking)
	admin.POST("/bookings/:id/reschedule", s.rescheduleBooking)
	admin.POST("/blocks", s.createBlock)
	admin.POST("/blocks/:id/delete", s.deleteBlock)
	admin.GET("/stats", s.stats)
	return r
}

func loadPages(base *template.Template) map[string]*template.Template {
	names, _ := fs.Glob(files, "templates/pages/*.html")
	pages := map[string]*template.Template{}
	for _, name := range names {
		pages[path.Base(name)] = template.Must(template.Must(base.Clone()).ParseFS(files, name))
	}
	return pages
}

var funcs = template.FuncMap{
	"initial": func(s string) string {
		r, _ := utf8.DecodeRuneInString(s)
		return string(r)
	},
	"pct": func(n, total int) int {
		if total == 0 {
			return 0
		}
		return n * 100 / total
	},
	"dict": func(kv ...any) map[string]any {
		m := map[string]any{}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		return m
	},
}

func (s *Server) render(c *gin.Context, page string, data any) {
	c.Render(http.StatusOK, render.HTML{Template: s.pages[page], Name: "layout", Data: data})
}

func (s *Server) partial(c *gin.Context, name string, data any) {
	c.Render(http.StatusOK, render.HTML{Template: s.partials, Name: name, Data: data})
}
