package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/damianpugliese/bookings/pkg/config"
	"github.com/damianpugliese/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewTeamplates(a *config.AppConfig)  {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}

	td.CSRFToken = nosurf.Token(r)

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache() 
	}

	// Get requested template from the cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	// Render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing teamplate to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from the ./templates/ folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	// Loop through the pages and create a template cache
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// Add the base layout to the template cache
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = ts
	}	

	return templateCache, nil
}