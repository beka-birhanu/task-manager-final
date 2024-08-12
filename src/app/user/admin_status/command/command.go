package promotcmd

import "github.com/google/uuid"

// Command represents a promote command with necessary data.
type Command struct {
	Username   string    // Username of the user to promote.
	PromoterID uuid.UUID // ID of the user performing the promotion.
}

// NewCommand creates a new Command instance with the given username and promoter ID.
func NewCommand(username string, promoterID uuid.UUID) *Command {
	return &Command{
		Username:   username,
		PromoterID: promoterID,
	}
}

