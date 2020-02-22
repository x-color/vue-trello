package model

// Board includes board data
type Board struct {
	ID     string
	UserID string
	Title  string
	Text   string
	Color  Color
}

// Boards defines a slice of Board
type Boards []Board
