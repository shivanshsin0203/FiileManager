package main

import (
	"filemanager/api"
	"filemanager/redis"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	r :=api.SetupRoutes()
    fmt.Println("Starting server on :8080...")
    redis.Init()
	defer redis.Close()
	corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),
        handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
    )(r)

	http.ListenAndServe(":8080", corsHandler)
}