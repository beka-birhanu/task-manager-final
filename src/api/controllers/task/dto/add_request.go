package dto

import (
	"time"
)

// DTOs for task operations
type AddTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}
