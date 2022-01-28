package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/and-cru/go-service/api/app/model"
	service "github.com/and-cru/go-service/api/app/service"
	"github.com/and-cru/go-service/api/config"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	// Create DB URI
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
	)

	// Connect to DB
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	fmt.Println("Ready")

	// auto migrate for local development, migrations for staging and prod
	if os.Getenv("ENV") == "develop" {
		a.DB = model.DBMigrate(db)
	} else {
		a.DB = db
	}

	// create and set router
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Health check for svc
	a.Get("/health", a.HealthCheck)

	// CRUD operation routing
	a.Get("/users", a.GetAllUsers)
	a.Post("/users", a.CreateUser)
	a.Get("/users/{name}", a.GetUser)
	a.Put("/users/{name}", a.UpdateUser)
	a.Delete("/users/{name}", a.DeleteUser)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	service.HealthChecker(w, r)
}

// services to manage Employee Data
func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	service.GetAllUsers(a.DB, w, r)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	service.CreateUser(a.DB, w, r)
}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	service.GetUser(a.DB, w, r)
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	service.UpdateUser(a.DB, w, r)
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	service.DeleteUser(a.DB, w, r)
}

func (a *App) DisableUser(w http.ResponseWriter, r *http.Request) {
	service.DisableUser(a.DB, w, r)
}

func (a *App) EnableUser(w http.ResponseWriter, r *http.Request) {
	service.EnableUser(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	//
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(NewOrigin("*"))
	n.UseHandler(a.Router)

	log.Fatal(http.ListenAndServe(host, n))
}
