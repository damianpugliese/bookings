package handlers

import (
	"net/http"

	"github.com/damianpugliese/bookings/pkg/config"
	"github.com/damianpugliese/bookings/pkg/models"
	"github.com/damianpugliese/bookings/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	romateIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", romateIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{
		"test": "Hello, again",
	}

	romateIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = romateIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}