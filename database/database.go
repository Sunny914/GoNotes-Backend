package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"gopkg.in/mgo.v2/bson"
)

// connect establishes connection to MongoDB and return the client and collections
func Connect() (*mongo.Client, *mongo.Collection, *mongo.Collection) {

	//connect to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB Connection String
	connectionString := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Collections
	db := client.Database("GoNotes")
	userCollection := db.Collection("users")
	noteCollection := db.Collection("notes")

	// Check Connection by Running a Query
	err = userCollection.FindOne(ctx, bson.M{}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal("Failed to query used collection: %v", err)
	}

	log.Println("Successfully Connected to the MongoDB database !")
	return client, userCollection, noteCollection

}





