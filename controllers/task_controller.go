package controllers

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/common"
	"github.com/beka-birhanu/controllers/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskController struct {
	taskService common.ITaskService
}

// NewTaskController creates a new TaskController with the given task service.
func NewTaskController(taskService common.ITaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// Register registers the task routes.
func (c *TaskController) Register(route gin.RouterGroup) {
	tasks := route.Group("/tasks")
	{
		tasks.POST("", c.addTask)
		tasks.PUT("/:id", c.updateTask)
		tasks.DELETE("/:id", c.deleteTask)
		tasks.GET("", c.getAllTasks)
		tasks.GET("/:id", c.getTask)
	}
}

func (c *TaskController) addTask(ctx *gin.Context) {
	var request dto.AddTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.taskService.Add(request.Title, request.Description, request.Status, request.DueDate)
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

func (c *TaskController) updateTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var request dto.TaskResponse

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = c.taskService.Update(id, request.Title, request.Description, request.Status, request.DueDate)
	if err != nil {
		if err == common.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *TaskController) deleteTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	err = c.taskService.Delete(id)
	if err != nil {
		if err == common.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *TaskController) getAllTasks(ctx *gin.Context) {
	tasks, err := c.taskService.GetAll()
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

func (c *TaskController) getTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := c.taskService.GetSingle(id)
	if err != nil {
		if err == common.ErrNotFound {
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
