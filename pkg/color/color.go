// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package color provides utility functions for formatting text with ANSI escape
// sequences to add terminal colors.
package color

import "fmt"

const escape = "\x1b"

const (
	NONE   = iota // NONE resets the terminal color formatting.
	RED           // RED sets the terminal text color to red.
	GREEN         // GREEN sets the terminal text color to green.
	YELLOW        // YELLOW sets the terminal text color to yellow.
	BLUE          // BLUE sets the terminal text color to blue.
	PURPLE        // PURPLE sets the terminal text color to purple.
)

// color generates the ANSI escape sequence string for the specified color code.
func color(c int) string {
	if c == NONE {
		return fmt.Sprintf("%s[%dm", escape, c)
	}

	return fmt.Sprintf("%s[3%dm", escape, c)
}

// Format wraps the provided text with ANSI escape sequences corresponding
// to the given color, resetting the color back to default at the end.
func Format(c int, text string) string {
	return color(c) + text + color(NONE)
}
