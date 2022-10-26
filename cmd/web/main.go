package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nk05081999/Reservation87/internal/config"
	"github.com/nk05081999/Reservation87/internal/driver"
	"github.com/nk05081999/Reservation87/internal/handlers"
	"github.com/nk05081999/Reservation87/internal/helpers"
	"github.com/nk05081999/Reservation87/internal/models"
	"github.com/nk05081999/Reservation87/internal/render"
)

// const portNumber = ":1010"

// func main() {
// 	var app config.AppConfig
// 	tc, err := render.CreateTemplateCache()
// 	if err != nil {
// 		log.Fatal("cannot create template cache")
// 	}
// 	app.TemplateCache = tc
// 	app.UseCache = false
// 	repo := handlers.NewRepo(&app)
// 	handlers.NewHandlers(repo)

// 	render.NewTemplates(&app)

// 	// http.HandleFunc("/", handlers.Repo.Home)
// 	 //http.HandleFunc("/about", handlers.Repo.About)

// 	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
// 	// _ = http.ListenAndServe(portNumber, nil)

// 	srv := &http.Server{
// 		Addr:    portNumber,
// 		Handler: routes(&app),
// 	}
// 	err = srv.ListenAndServe()
// 	// if err != nil {
// 	log.Fatal(err)
// 	// }
// }

const portNumber = ":1010"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func run() (*driver.DB, error) {

	// what i am going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.Errorlog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to the database....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=practice user=postgres password=root")
	if err != nil {
		log.Fatal("cannot connect to the database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
