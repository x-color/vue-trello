package model

// Board includes board data
type Board struct {
	ID     string
	UserID string
	Title  string
	Text   string
	Color  Color
	Lists  Lists
}

// Boards defines a slice of Board
type Boards []Board
