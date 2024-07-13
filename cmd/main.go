package main

import (
	"fifthtask/internal/db"
	"fifthtask/internal/handlers"
	"log"
	"net/http"
)

const conn string = "localhost:8080"

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/profile", handlers.ProfileHandler)

	log.Println("Server running on", conn)
	log.Fatal(http.ListenAndServe(conn, mux))
}
