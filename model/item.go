package model

// Item includes item data
type Item struct {
	ID     string
	ListID string
	UserID string
	Title  string
	Text   string
	Tags   Tags
}

// Items defines a slice of Item
type Items []Item
