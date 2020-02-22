package model

// User includes user data
type User struct {
	ID   string
	Name string
	// Password is raw. It will be hash.
	Password string
}
