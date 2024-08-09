package loginqry

// Query represents the data required for a login request.
type Query struct {
	Username string // Username of the user trying to log in.
	Password string // Password of the user trying to log in.
}

// NewQuery creates a new Query instance with the provided username and password.
func NewQuery(username, password string) *Query {
	return &Query{
		Username: username,
		Password: password,
	}
}
