package err

// Validation errors
var (
	ErrTitleEmpty       = NewValidation("title cannot be empty")
	ErrDescriptionEmpty = NewValidation("description cannot be empty")
	ErrDueDateZero      = NewValidation("due date cannot be zero")
	ErrInvalidStatus    = NewValidation("invalid status")
)
