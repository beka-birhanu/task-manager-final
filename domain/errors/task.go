package err

// Validation errors
var (
	// ErrTitleEmpty indicates that the title cannot be empty.
	ErrTitleEmpty = NewValidation("title cannot be empty")

	// ErrDescriptionEmpty indicates that the description cannot be empty.
	ErrDescriptionEmpty = NewValidation("description cannot be empty")

	// ErrDueDateZero indicates that the due date cannot be zero.
	ErrDueDateZero = NewValidation("due date cannot be zero")

	// ErrInvalidStatus indicates that the status is invalid.
	ErrInvalidStatus = NewValidation("invalid status")
)

