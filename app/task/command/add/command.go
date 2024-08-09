package addcmd

import "time"

// Command represents the data required to add a new task.
// Fields:
// - title: The title of the task.
// - description: A detailed description of the task.
// - status: The current status of the task.
// - dueDate: The due date for the task.
type Command struct {
	title       string
	description string
	status      string
	dueDate     time.Time
}

// NewCommand creates a new Command instance with the specified details.
func NewCommand(title, description, status string, dueDate time.Time) *Command {
	return &Command{
		title:       title,
		description: description,
		status:      status,
		dueDate:     dueDate,
	}
}

