package errdmn

// Validation errors
var (
	// TitleEmpty indicates that the title cannot be empty.
	TitleEmpty = NewValidation("title cannot be empty")

	// DescriptionEmpty indicates that the description cannot be empty.
	DescriptionEmpty = NewValidation("description cannot be empty")

	// DueDateZero indicates that the due date cannot be zero.
	DueDateZero = NewValidation("due date cannot be zero")

	// InvalidStatus indicates that the status is invalid.
	InvalidStatus = NewValidation("invalid status")

	// TaskNotFound indicates that a task was not found.
	TaskNotFound = NewValidation("task not found")
)

