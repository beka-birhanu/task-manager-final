package taskcontroller

import (
	"fmt"
	"net/http"

	basecontroller "github.com/beka-birhanu/task_manager_final/api/controllers/base"
	"github.com/beka-birhanu/task_manager_final/api/controllers/task/dto"
	errapi "github.com/beka-birhanu/task_manager_final/api/errors"
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	addcmd "github.com/beka-birhanu/task_manager_final/app/task/command/add"
	updatecmd "github.com/beka-birhanu/task_manager_final/app/task/command/update"
	errdmn "github.com/beka-birhanu/task_manager_final/domain/errors"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller handles HTTP requests related to tasks.
type Controller struct {
	basecontroller.BaseHandler
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

// RegisterPublic registers public routes.
func (c *Controller) RegisterPublic(route *gin.RouterGroup) {}

// RegisterProtected registers protected routes.
func (c *Controller) RegisterProtected(route *gin.RouterGroup) {
	tasks := route.Group("/tasks")
	{
		tasks.GET("", c.getAllTasks)
		tasks.GET("/:id", c.getTask)
	}
}

// RegisterPrivileged registers privileged routes.
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
		c.Problem(ctx, errapi.NewBadRequest(err.Error()))
		return
	}

	cmd := addcmd.NewCommand(request.Title, request.Description, request.Status, request.DueDate)
	task, err := c.addHandler.Handle(cmd)
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		DueDate:     task.DueDate(),
		Status:      task.Status(),
	}

	baseURL := fmt.Sprintf("http://%s", ctx.Request.Host)
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, ctx.Request.URL.Path, task.ID().String())
	c.RespondWithLocation(ctx, http.StatusCreated, response, resourceLocation)
}

func (c *Controller) updateTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	var request dto.AddTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.Problem(ctx, errapi.NewBadRequest(err.Error()))
		return
	}

	cmd := updatecmd.NewCommand(id, request.Title, request.Description, request.Status, request.DueDate)
	_, err = c.updateHandler.Handle(cmd)
	if err != nil {
		if err == errdmn.TaskNotFound {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		} else {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		}
		return
	}

	c.Respond(ctx, http.StatusOK, nil)
}

func (c *Controller) deleteTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	_, err = c.deleteHandler.Handle(id)
	if err != nil {
		if err == errdmn.TaskNotFound {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		} else {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		}
		return
	}

	c.Respond(ctx, http.StatusOK, nil)
}

func (c *Controller) getAllTasks(ctx *gin.Context) {
	tasks, err := c.getAllHandler.Handle(struct{}{})
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
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
	c.Respond(ctx, http.StatusOK, response)
}

func (c *Controller) getTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Problem(ctx, errapi.NewBadRequest(err.Error()))
		return
	}

	task, err := c.getHandler.Handle(id)
	if err != nil {
		if err == errdmn.TaskNotFound {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		} else {
			c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
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

	c.Respond(ctx, http.StatusOK, response)
}
