package main

import (
	"log"
	"net/http"
	"os"

	"yal/app"
	"yal/controllers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// NOTE: auth required to see registered users
	router.HandleFunc("/api/users",
		controllers.GetUsers).Methods("GET")
	router.HandleFunc("/api/user/{id}",
		controllers.GetUser).Methods("GET")

	// no auth for register and login
	router.HandleFunc("/api/user/new",
		controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login",
		controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/user/{id}/update",
		controllers.UpdateAccount).Methods("PUT")

	// JWT middleware
	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// INFO: api will be located at `localhost:8080/api/...`
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalln(err)
	}
}
