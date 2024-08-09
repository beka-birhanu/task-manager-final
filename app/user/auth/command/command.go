package registercmd

// Command represents the data required for user registration.
type Command struct {
	Username string // Username of the new user.
	Password string // Password for the new user.
}

// NewCommand creates a new Command instance with the specified username and password.
func NewCommand(username, password string) *Command {
	return &Command{
		Username: username,
		Password: password,
	}
}

