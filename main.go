package main

import (
	"context"
	"log"

	taskcontrollers "github.com/beka-birhanu/controllers/task"
	taskrepo "github.com/beka-birhanu/data/task"
	"github.com/beka-birhanu/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration constants
const (
	addr    = ":8080"
	baseURL = "/api"
)

func main() {
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI("")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}

	// Ping the MongoDB server to verify connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("Error pinging MongoDB server: %v", err)
	}

	// Create a new task service instance
	taskService := taskrepo.New(client, "taskdb", "tasks")

	// Create a new task controller instance
	taskController := taskcontrollers.New(taskService)

	// Create a new router instance with configuration
	routerConfig := router.Config{
		Addr:         addr,
		BaseURL:      baseURL,
		TasksHandler: taskController,
	}
	router := router.NewRouter(routerConfig)

	// Run the server
	if err := router.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
