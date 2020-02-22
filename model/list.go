package model

// List includes list data
type List struct {
	ID      string
	BoardID string
	UserID  string
	Title   string
}

// Lists defines a slice of List
type Lists []List
