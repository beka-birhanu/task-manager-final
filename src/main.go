package main

import (
	"fmt"
	"log"

	"github.com/beka-birhanu/task_manager_final/src/api"
	authcontroller "github.com/beka-birhanu/task_manager_final/src/api/controllers/auth"
	taskcontroller "github.com/beka-birhanu/task_manager_final/src/api/controllers/task"
	usercontroller "github.com/beka-birhanu/task_manager_final/src/api/controllers/user"
	"github.com/beka-birhanu/task_manager_final/src/api/router"
	addcmd "github.com/beka-birhanu/task_manager_final/src/app/task/command/add"
	deletecmd "github.com/beka-birhanu/task_manager_final/src/app/task/command/delete"
	updatecmd "github.com/beka-birhanu/task_manager_final/src/app/task/command/update"
	getqry "github.com/beka-birhanu/task_manager_final/src/app/task/query/get"
	getallqry "github.com/beka-birhanu/task_manager_final/src/app/task/query/get_all"
	promotcmd "github.com/beka-birhanu/task_manager_final/src/app/user/admin_status/command"
	registercmd "github.com/beka-birhanu/task_manager_final/src/app/user/auth/command"
	loginqry "github.com/beka-birhanu/task_manager_final/src/app/user/auth/query"
	"github.com/beka-birhanu/task_manager_final/src/config"
	"github.com/beka-birhanu/task_manager_final/src/infrastructure/db"
	"github.com/beka-birhanu/task_manager_final/src/infrastructure/hash"
	"github.com/beka-birhanu/task_manager_final/src/infrastructure/jwt"
	taskrepo "github.com/beka-birhanu/task_manager_final/src/infrastructure/repo/task"
	userrepo "github.com/beka-birhanu/task_manager_final/src/infrastructure/repo/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// main is the entry point for the application.
// It initializes the MongoDB client, services, controllers, and starts the HTTP server.
func main() {
	cfg := config.Envs

	// Initialize MongoDB client and perform migrations
	mongoClient := initDB(cfg)

	// Initialize services
	userRepo, taskRepo, jwtService, hashService := initServices(cfg, mongoClient)

	// Initialize controllers
	userController := initUserController(userRepo)
	authController := initAuthController(userRepo, jwtService, hashService)
	taskController := initTaskController(taskRepo)

	// Router configuration
	routerConfig := router.Config{
		Addr:        fmt.Sprintf(":%s", cfg.ServerPort),
		BaseURL:     "/api",
		Controllers: []api.IController{userController, taskController, authController},
		JwtService:  jwtService,
	}
	r := router.NewRouter(routerConfig)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// initDB initializes the MongoDB client and performs any necessary database migrations.
// It returns the MongoDB client instance.
func initDB(cfg config.Config) *mongo.Client {
	mongoClient := db.Connect(db.Config{
		ConnectionString: cfg.DBConnectionString,
	})

	db.Migrate(mongoClient, cfg.DBName)

	return mongoClient
}

// initServices initializes the necessary services for the application.
// It returns the user repository, task repository, JWT service, and hash service.
func initServices(cfg config.Config, mongoClient *mongo.Client) (*userrepo.Repo, *taskrepo.Repo, *jwt.Service, *hash.Service) {
	userRepo := userrepo.NewRepo(mongoClient, cfg.DBName, "users")
	taskRepo := taskrepo.New(mongoClient, cfg.DBName, "tasks")

	jwtService := jwt.New(jwt.Config{
		SecretKey: cfg.JWTSecret,
		Issuer:    cfg.ServerHost,
		ExpTime:   cfg.JWTExpirationInSeconds,
	})

	hashService := hash.SingletonService()

	return userRepo, taskRepo, jwtService, hashService
}

// initUserController initializes the user controller with the necessary handlers.
// It returns the user controller instance.
func initUserController(userRepo *userrepo.Repo) *usercontroller.Controller {
	promotHandler := promotcmd.New(userRepo)

	return usercontroller.New(usercontroller.Config{
		PromotHandler: promotHandler,
	})
}

// initAuthController initializes the authentication controller with the necessary handlers.
// It returns the authentication controller instance.
func initAuthController(userRepo *userrepo.Repo, jwtService *jwt.Service, hashService *hash.Service) *authcontroller.Controller {
	signupHandler := registercmd.NewHandler(registercmd.Config{
		UserRepo: userRepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
	})

	loginHandler := loginqry.NewHandler(loginqry.Config{
		UserRepo: userRepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
	})

	return authcontroller.New(authcontroller.Config{
		RegisterHandler: signupHandler,
		LoginHandler:    loginHandler,
	})
}

// initTaskController initializes the task controller with the necessary handlers.
// It returns the task controller instance.
func initTaskController(taskRepo *taskrepo.Repo) *taskcontroller.Controller {
	addHandler := addcmd.NewHandler(taskRepo)
	updateHandler := updatecmd.NewHandler(taskRepo)
	deleteHandler := deletecmd.New(taskRepo)
	getAllHandler := getallqry.New(taskRepo)
	getHandler := getqry.New(taskRepo)

	return taskcontroller.New(taskcontroller.Config{
		AddHandler:    addHandler,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
		GetAllHandler: getAllHandler,
		GetHandler:    getHandler,
	})
}
