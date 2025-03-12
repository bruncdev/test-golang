package main

import (
	"log"
	"net/http"
	"os"

	"people-api/database"
	"people-api/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	database.InitDB()
	defer database.DB.Close()

	// Create router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/people", handlers.CreatePerson).Methods("POST")
	r.HandleFunc("/api/people", handlers.GetPeople).Methods("GET")
	r.HandleFunc("/api/people/{id}", handlers.GetPerson).Methods("GET")
	r.HandleFunc("/api/people/{id}", handlers.UpdatePerson).Methods("PUT")
	r.HandleFunc("/api/people/{id}", handlers.DeletePerson).Methods("DELETE")

	// Set up middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
