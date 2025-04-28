package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/damianpugliese/bookings/pkg/config"
	"github.com/damianpugliese/bookings/pkg/handlers"
	"github.com/damianpugliese/bookings/pkg/render"
)

const port = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTeamplates(&app)

	fmt.Println("Starting server on port", port)

	svr := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = svr.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}