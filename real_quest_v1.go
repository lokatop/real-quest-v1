package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"real-quest-v1/app"
	"real-quest-v1/controllers"
	u "real-quest-v1/models"
	"time"
)


func main() {
	// Initialize a connection pool and assign it to the pool global
	// variable.
	u.Pool = &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/quest/{id:[0-9]+}", controllers.ShowQuest).Methods("GET")
	router.HandleFunc("/like", controllers.AddLike).Methods("POST")
	router.HandleFunc("/popular", controllers.ListPopular).Methods("GET")
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	//router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")

	router.HandleFunc("/new", controllers.CreateRecord).Methods("GET")

	//down ↓ var NotAuth uses for definition router without authentication
	app.NotAuth = []string{"/like","/new", "/api/user/new", "/api/user/login","/popular","/quest/{id:[0-9]+}"}

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}