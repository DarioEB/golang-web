package main

import (
	"encoding/json"
	"fmt"
	"go_sample_api/internal/course"
	"go_sample_api/internal/user"
	"go_sample_api/pkg/bootstrap"
	"log"
	"net/http"
	"os"

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
	courseRepo := course.NewRepo(logger, db)
	userSvc := user.InitService(logger, userRepo)
	courseSvc := course.NewService(logger, courseRepo)
	userEnd := user.MakeEndpoints(logger, userSvc)
	courseEnd := course.MakeEndpoints(logger, courseSvc)
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PUT")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")
	router.HandleFunc("/course", courseEnd.Create).Methods("POST")

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT")),
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
