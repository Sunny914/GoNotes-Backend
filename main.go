package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Sunny914/GoNotes-Backend/database"
	"github.com/Sunny914/GoNotes-Backend/handlers"
	"github.com/Sunny914/GoNotes-Backend/models"
	"github.com/Sunny914/GoNotes-Backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (likely running in production):", err)
	}

	// Connect to MongoDB
	client, userCollection, noteCollection := database.Connect()
	defer client.Disconnect(context.Background())

	// Create Models
	userModel := models.NewUserModel(userCollection)
	noteModel := models.NewNoteModel(noteCollection, userCollection)

	// Define your JWT secret key (keep it safe and strong)
	jwtSecret := []byte("your-secret-key") // Replace with a secure secret

	// Create handlers with JWT-based auth
	authHandler := handlers.NewAuthHandler(userModel, jwtSecret)
	noteHandler := handlers.NewNoteHandler(noteModel, jwtSecret)

	// configure router
	r := mux.NewRouter()
	routes.Setup(r, authHandler, noteHandler)

	// start server
	log.Println("Server starting at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}