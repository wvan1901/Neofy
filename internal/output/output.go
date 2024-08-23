package output

import (
	"fmt"
	"neofy/internal/config"
)

// This is the main draw func
func UpdateApp(d *config.AppData) {
	d.Display.Buffer.WriteString("\033[2J") // Clears entire screen
	//d.Display.Buffer.WriteString("\033[K") // Clears entire Line
	// Clears Screen
	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor
	d.Display.Buffer.WriteString("\033[H")    // Move Cursor to upper right
	// Draw everything
	d.Display.Buffer.WriteString("Wilson Was Here!")

	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor

	fmt.Print(d.Display.Buffer.String())
	d.Display.Buffer.Reset()
}
