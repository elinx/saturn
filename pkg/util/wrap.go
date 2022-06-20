package util

import (
	"unicode/utf8"

	"github.com/zyedidia/go-runewidth"
)

func IsTerminator(c rune) bool {
	return (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a)
}

// Wrap wraps the given string to the given width.
// The result is a string with newlines inserted at the right places.
// The result is guaranteed to be shorter than the given width.
// Newlines in the input string are preserved.
// Leading and trailing whitespaces are preserved.
// Ansi escape sequences are preserved.
// Tab keys are not expanded.
func Wrap(line string, limit int) string {
	if len(line) <= limit || limit < 2 {
		return line
	}
	lineWidth := 0
	result := ""
	ansi := false
	for len(line) > 0 {
		rune, size := utf8.DecodeRuneInString(line)
		if rune == '\x1b' {
			ansi = true
			result += string(rune)
			line = line[size:]
			continue
		}
		if ansi {
			if IsTerminator(rune) {
				ansi = false
			}
			result += string(rune)
			line = line[size:]
			continue
		}
		if rune == '\n' {
			lineWidth = 0
			line = line[size:]
			result += "\n"
			continue
		}
		// FIXME(elinx): support tab key
		cellWidth := runewidth.RuneWidth(rune)
		if lineWidth+cellWidth <= limit {
			result += string(rune)
			lineWidth += cellWidth
			line = line[size:]
		} else if lineWidth+cellWidth > limit {
			result += "\n"
			lineWidth = 0
		}
	}
	return result
}

// LocBeforeWraped returns the original position of the character before the
// string is wraped given the width. vx and vy are the visual position of the
// character. The result is the rune index of the character before the wraped.
func LocBeforeWraped(line string, limit int, vx, vy int) int {
	if len(line) <= limit || vy == 0 {
		return vx
	}
	ansi := false
	lineWidth := 0
	x := 0
	for len(line) > 0 {
		rune, size := utf8.DecodeRuneInString(line)
		if rune == '\x1b' {
			ansi = true
			line = line[size:]
			continue
		}
		if ansi {
			if IsTerminator(rune) {
				ansi = false
			}
			line = line[size:]
			continue
		}
		if rune == '\n' {
			lineWidth = 0
			line = line[size:]
			x++
			vy--
			continue
		}
		cellWidth := runewidth.RuneWidth(rune)
		if lineWidth+cellWidth <= limit {
			lineWidth += cellWidth
			line = line[size:]
			x++
		} else if lineWidth+cellWidth > limit {
			lineWidth = 0
			vy--
		}
		if vy == 0 && vx <= lineWidth {
			return x
		}
	}
	return x
}
