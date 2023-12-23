package main

import (
	"encoding/json"
	"fmt"
	"go_sample_api/internal/user"
	"go_sample_api/pkg/bootstrap"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	_ = godotenv.Load()
	logger := bootstrap.InitLogger()
	db, err := bootstrap.Connection(logger)
	if err != nil {
		logger.Fatal(err)
	}
	userRepo := user.NewRepo(logger, db)
	userSvc := user.InitService(logger, userRepo)
	userEnd := user.MakeEndpoints(logger, userSvc)
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PUT")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3000",
	}

	log.Fatal(srv.ListenAndServe())
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /users")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /courses")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
