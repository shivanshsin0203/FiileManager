package api

import (
	"filemanager/auth"
	"filemanager/aws"
	"filemanager/redis"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/validate", auth.Validate).Methods("GET")
	router.HandleFunc(("/generatePresignedURL"), aws.GeneratePresignedURL).Methods("GET")
	router.HandleFunc(("/addQueue"),rediss.EnqueueHandler).Methods("POST")
	return router
}