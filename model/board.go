package model

// Board includes board data
type Board struct {
	ID     string
	UserID string
	Title  string
	Text   string
	Color  Color
	Lists  Lists
	Before string
	After  string
}

// Boards defines a slice of Board
type Boards []Board
