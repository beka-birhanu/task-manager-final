package taskcontroller

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/api/controllers/task/dto"
	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	addcmd "github.com/beka-birhanu/app/task/command/add"
	updatecmd "github.com/beka-birhanu/app/task/command/update"
	errdmn "github.com/beka-birhanu/domain/errors"
	taskmodel "github.com/beka-birhanu/domain/models/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller handles HTTP requests related to tasks.
type Controller struct {
	addHandler    icmd.IHandler[*addcmd.Command, *taskmodel.Task]
	updateHandler icmd.IHandler[*updatecmd.Command, *taskmodel.Task]
	deleteHandler icmd.IHandler[uuid.UUID, bool]
	getAllHandler icmd.IHandler[struct{}, []*taskmodel.Task]
	getHandler    icmd.IHandler[uuid.UUID, *taskmodel.Task]
}

type Config struct {
	AddHandler    icmd.IHandler[*addcmd.Command, *taskmodel.Task]
	UpdateHandler icmd.IHandler[*updatecmd.Command, *taskmodel.Task]
	DeleteHandler icmd.IHandler[uuid.UUID, bool]
	GetAllHandler icmd.IHandler[struct{}, []*taskmodel.Task]
	GetHandler    icmd.IHandler[uuid.UUID, *taskmodel.Task]
}

// New creates a new TaskController with the given CQRS handlers and task repository.
func New(config Config) *Controller {
	return &Controller{
		addHandler:    config.AddHandler,
		updateHandler: config.UpdateHandler,
		deleteHandler: config.DeleteHandler,
		getAllHandler: config.GetAllHandler,
		getHandler:    config.GetHandler,
	}
}

// Register registers the task routes with CQRS handlers.
func (c *Controller) Register(route *gin.RouterGroup) {

}
func (c *Controller) RegisterPublic(route *gin.RouterGroup) {}

func (c *Controller) RegisterProtected(route *gin.RouterGroup) {
	tasks := route.Group("/tasks")
	{
		tasks.GET("", c.getAllTasks)
		tasks.GET("/:id", c.getTask)
	}

}

func (c *Controller) RegisterPrivileged(route *gin.RouterGroup) {
	tasks := route.Group("/tasks")
	{
		tasks.POST("", c.addTask)
		tasks.PUT("/:id", c.updateTask)
		tasks.DELETE("/:id", c.deleteTask)

	}
}

func (c *Controller) addTask(ctx *gin.Context) {
	var request dto.AddTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := addcmd.NewCommand(request.Title, request.Description, request.Status, request.DueDate)
	task, err := c.addHandler.Handle(cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the base URL
	baseURL := fmt.Sprintf("http://%s", ctx.Request.Host)

	// Construct the resource location
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, ctx.Request.URL.Path, task.ID().String())

	// Set the Location header and return the response
	ctx.Header("Location", resourceLocation)
	ctx.Status(http.StatusCreated)
}

func (c *Controller) updateTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var request dto.AddTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := updatecmd.NewCommand(id, request.Title, request.Description, request.Status, request.DueDate)
	_, err = c.updateHandler.Handle(cmd)
	if err != nil {
		if err == errdmn.TaskNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) deleteTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	_, err = c.deleteHandler.Handle(id)
	if err != nil {
		if err == errdmn.TaskNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) getAllTasks(ctx *gin.Context) {
	tasks, err := c.getAllHandler.Handle(struct{}{})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var response []dto.TaskResponse
	for _, task := range tasks {
		response = append(response, dto.TaskResponse{
			ID:          task.ID(),
			Title:       task.Title(),
			Description: task.Description(),
			DueDate:     task.DueDate(),
			Status:      task.Status(),
		})
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) getTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := c.getHandler.Handle(id)
	if err != nil {
		if err == errdmn.TaskNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		DueDate:     task.DueDate(),
		Status:      task.Status(),
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
