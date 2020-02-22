package model

// Tag includes tags data
type Tag struct {
	ID    string
	Name  string
	Color Color
}

// Tags defines a slice of Tag
type Tags []Tag
