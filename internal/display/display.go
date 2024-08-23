package display

import (
	"strings"
)

type Display struct {
	Buffer strings.Builder
	Width  int
	Height int
}

func InitDisplay(w, h int) *Display {
	var newBuf strings.Builder
	newBuf.Reset()
	display := Display{
		Buffer: newBuf,
		Height: h,
		Width:  w,
	}
	return &display
}

func (d *Display) WriteString(s string) {
	d.Buffer.WriteString(s)
}

func clearScreen() {}

func drawScreen() {
	clearScreen()
}
