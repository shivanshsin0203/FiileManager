package api

import (
	"filemanager/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/validate", auth.Validate).Methods("GET")
	return router
}