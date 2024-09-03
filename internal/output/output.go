package output

import (
	"fmt"
	"neofy/internal/config"
	"strconv"
	"strings"
)

// This is the main draw func
func UpdateApp(d *config.AppData) {
	d.Display.Buffer.WriteString("\033[2J") // Clears entire screen
	//d.Display.Buffer.WriteString("\033[K") // Clears entire Line
	// Clears Screen
	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor
	d.Display.Buffer.WriteString("\033[H")    // Move Cursor to upper right

	// Update App Components
	updatePlayerDisplay(&d.Player)

	// Draw everything
	drawScreen(50, 25, d)

	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor

	fmt.Print(d.Display.Buffer.String())
	d.Display.Buffer.Reset()
}

func drawScreen(width, height int, d *config.AppData) {
	d.Display.Buffer.WriteString("Neofy v0.0.0 | " + "w:" + strconv.Itoa(width) + " h:" + strconv.Itoa(height) + "\r\n")
	printPlayerView(d.Player.Display.Screen, &d.Display.Buffer)
}

func printPlayerView(s []string, buf *strings.Builder) {
	for i := range s {
		buf.WriteString(s[i])
		buf.WriteString("\r\n")
	}
}

// TODO: If len of row is too long trim it
func updatePlayerDisplay(mp *config.MusicPlayer) {
	// If there is no screen then we do nothing
	s := mp.Display.Screen
	if len(s) == 0 {
		return
	}

	// Prepare visual for player
	playPause := "|>"
	if mp.IsPlaying {
		playPause = "||"
	}
	volume := "<|)"
	if mp.SupportsVolume {
		volume += " "
		numBold := mp.Volume / 10
		if numBold > 10 {
			numBold = 10
		} else if numBold < 0 {
			numBold = 0
		}
		numSlim := 10 - numBold
		bold := strings.Repeat("=", numBold)
		slim := strings.Repeat("-", numSlim)
		volume += bold + slim
	} else {
		volume += " xxxxx"
	}
	shuffled := "->"
	if mp.IsShuffled {
		shuffled = "x>"
	}
	loop := "???"
	switch mp.Repeat {
	case "off":
		loop = "~~>"
	case "context":
		loop = "[≥]"
	case "track":
		loop = "[!]"
	}

	// Write Visual to display
	for i := range s {
		switch i {
		case len(s) / 2:
			// TODO: Make the player appear in the middle
			//"x>    |<    ||    >|    [≥]    <|) =====-----"
			s[i] = shuffled + "    |<    " + playPause + "    >|    " + loop + "    " + volume
		case 0:
			if mp.CurrentSong.Name == "" {
				s[i] = ""
			} else {
				s[i] = "SongPlaying: " + mp.CurrentSong.Name
			}
		case len(s) - 1:
			if mp.CurrentSong.Artist == "" {
				s[i] = ""
			} else {
				s[i] = "SongArtist: " + mp.CurrentSong.Artist
			}
		}
	}
}
