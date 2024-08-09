package main

import (
	"context"
	"log"

	taskcontrollers "github.com/beka-birhanu/controllers/task"
	usercontroller "github.com/beka-birhanu/controllers/user"
	taskrepo "github.com/beka-birhanu/infrastructure/repo/task"
	userrepo "github.com/beka-birhanu/infrastructure/repo/user"
	"github.com/beka-birhanu/router"
	"github.com/beka-birhanu/service/hash"
	"github.com/beka-birhanu/service/jwt"
	usersvc "github.com/beka-birhanu/service/user"
	"go.mongodb.org/mongo-driver/bson"
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
	// Choose the database and collection
	database := client.Database("taskdb")
	usersCollection := database.Collection("users")

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"username": 1, // 1 for ascending order
		},
		Options: options.Index().SetUnique(true), // Optional: Set the index to be unique
	}

	// Create the index
	indexName, err := usersCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created index %s\n", indexName)
	// Create a new task service instance
	taskService := taskrepo.New(client, "taskdb", "tasks")

	// Create a new task controller instance
	taskController := taskcontrollers.New(taskService)

	jwtService := jwt.New(jwt.Config{
		SecretKey: "not-so-secret-now-is-it?",
		ExpTime:   1400,
	})

	hashService := hash.SingletonService()
	userrepo := userrepo.NewUserRepo(client, "taskdb", "usres")
	usersvc := usersvc.NewService(usersvc.Config{
		UserRepo: userrepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
	})
	usercontroller := usercontroller.New(usersvc)

	// Create a new router instance with configuration
	routerConfig := router.Config{
		Addr:         addr,
		BaseURL:      baseURL,
		TasksHandler: taskController,
		UserHandler:  usercontroller,
	}
	router := router.NewRouter(routerConfig)

	// Run the server
	if err := router.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
