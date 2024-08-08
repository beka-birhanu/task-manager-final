package router

import (
	"log"

	"github.com/beka-birhanu/common"
	"github.com/gin-gonic/gin"
)

// Router is the struct for managing the HTTP server and its dependencies.
type Router struct {
	addr         string
	baseURL      string
	tasksHandler common.IController
}

// Config is the struct for configuring the Router.
type Config struct {
	Addr         string
	BaseURL      string
	TasksHandler common.IController
}

// NewRouter creates a new instance of Router with the given configuration.
func NewRouter(config Config) *Router {
	return &Router{
		addr:         config.Addr,
		baseURL:      config.BaseURL,
		tasksHandler: config.TasksHandler,
	}
}

// Run starts the HTTP server and sets up the routes.
func (r *Router) Run() error {
	router := gin.Default()

	// Setting up routes under baseURL
	api := router.Group(r.baseURL)
	{
		v1 := api.Group("/v1")
		{
			r.tasksHandler.Register(*v1)
		}
	}

	log.Println("Listening on", r.addr)
	return router.Run(r.addr)
}
