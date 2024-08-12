package updatecmd

import (
	"time"

	"github.com/google/uuid"
)

// Command represents the data needed to update an existing task.
type Command struct {
	id          uuid.UUID
	title       string
	description string
	status      string
	dueDate     time.Time
}

// NewCommand creates a new Command instance with the provided task details.
func NewCommand(id uuid.UUID, title, description, status string, dueDate time.Time) *Command {
	return &Command{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		dueDate:     dueDate,
	}
}
