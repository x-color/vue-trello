package model

// Color defines a color type
type Color string

// Color pattern
const (
	RED    Color = "red"
	BLUE   Color = "blue"
	GREEN  Color = "green"
	YELLOW Color = "yellow"
)

// Colors defines a slice of Color.
type Colors []Color

// COLORS includes Color literals.
var COLORS = Colors{
	RED,
	BLUE,
	GREEN,
	YELLOW,
}
