package web

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path"

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
}

type Server struct {
	store Store
	pages map[string]*template.Template
}

func New(st Store) *Server {
	return &Server{store: st, pages: loadPages()}
}

func (s *Server) Routes() *gin.Engine {
	r := gin.Default()
	static, _ := fs.Sub(files, "static")
	r.StaticFS("/static", http.FS(static))
	r.GET("/", s.index)
	r.POST("/booking", s.createBooking)
	return r
}

func loadPages() map[string]*template.Template {
	base := template.Must(template.ParseFS(files, "templates/layout.html"))
	names, _ := fs.Glob(files, "templates/pages/*.html")
	pages := map[string]*template.Template{}
	for _, name := range names {
		pages[path.Base(name)] = template.Must(template.Must(base.Clone()).ParseFS(files, name))
	}
	return pages
}

func (s *Server) render(c *gin.Context, page string, data any) {
	c.Render(http.StatusOK, render.HTML{Template: s.pages[page], Name: "layout", Data: data})
}
