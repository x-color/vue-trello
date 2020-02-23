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

// COLORS includes Color literals.
const COLORS = []Color{
	RED,
	BLUE,
	GREEN,
	YELLOW,
}
