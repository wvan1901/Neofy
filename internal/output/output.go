package output

import (
	"fmt"
	"neofy/internal/data"
	"strings"
	"unicode/utf8"
)

// This is the main draw func
func UpdateApp(d *data.AppData) {
	d.Display.Buffer.WriteString("\033[2J") // Clears entire screen
	//d.Display.Buffer.WriteString("\033[K") // Clears entire Line
	// Clears Screen
	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor
	d.Display.Buffer.WriteString("\033[H")    // Move Cursor to upper right

	// Update App Components
	updatePlaylistDisplay(&d.Playlist)
	updateTracksDisplay(&d.Songs)
	updatePlayerDisplay(&d.Player)

	// Draw everything
	drawAppScreen(d)

	d.Display.Buffer.WriteString("\033[?25l") // Moves Cursor

	fmt.Print(d.Display.Buffer.String())
	d.Display.Buffer.Reset()
}

func drawAppScreen(d *data.AppData) {
	d.Display.Buffer.WriteString("Neofy v0.0.0\r\n")
	drawMusicOptions(&d.Playlist, &d.Songs, &d.Display.Buffer) // Playlist & tracks
	drawPlayer(&d.Player, &d.Display.Buffer)

}

func drawMusicOptions(playlist *data.Playlist, tracks *data.Tracks, buf *strings.Builder) {
	if playlist.Display.Height != tracks.Display.Height {
		return
	}
	numRows := playlist.Display.Height
	for i := 0; i < numRows; i++ {
		rowString := playlist.Display.Screen[i] + "|" + tracks.Display.Screen[i] + "|" + "\r\n"
		buf.WriteString(rowString)
	}
}

func drawPlayer(p *data.MusicPlayer, buf *strings.Builder) {
	if p.Display.Height < 3 {
		return
	}
	for i := range p.Display.Screen {
		rowString := p.Display.Screen[i] + "|\r\n"
		//rowString := p.Display.Screen[i] + "\r\n"
		buf.WriteString(rowString)
	}
}

func printPlayerView(s []string, buf *strings.Builder) {
	for i := range s {
		buf.WriteString(s[i])
		buf.WriteString("\r\n")
	}
}

// TODO: If len of row is too long trim it
func updatePlayerDisplay(mp *data.MusicPlayer) {
	// If there is no screen then we do nothing
	s := []string{}

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
	header := fitStringToWidthAndFillRune("--Song", '-', mp.Display.Width)
	s = append(s, header)
	for i := 0; i < mp.Display.Height-2; i++ {
		str := fitStringToWidth("&", mp.Display.Width)
		switch i {
		case (mp.Display.Height - 2) / 2:
			// TODO: Make the player appear in the middle
			//"x>    |<    ||    >|    [≥]    <|) =====-----"
			str = fitStringInMiddle(shuffled+"    |<    "+playPause+"    >|    "+loop+"    "+volume, ' ', mp.Display.Width)
		case 0:
			if mp.CurrentSong.Name == "" {
				str = fitStringToWidth("", mp.Display.Width)
			} else {
				str = fitStringToWidth("Playing: "+mp.CurrentSong.Name, mp.Display.Width)
			}
		case mp.Display.Height - 3:
			if mp.CurrentSong.Artist == "" {
				str = fitStringToWidth("", mp.Display.Width)
			} else {
				str = fitStringToWidth("Artist: "+mp.CurrentSong.Artist, mp.Display.Width)
			}
		}
		s = append(s, str)
	}
	bottom := fillWidthWithRune('-', mp.Display.Width)
	s = append(s, bottom)
	mp.Display.Screen = s
}

func updatePlaylistDisplay(playlist *data.Playlist) {
	// TODO: Handle offet
	playlist.Display.Screen = []string{}
	header := fitStringInMiddle("Playlists", '-', playlist.Display.Width)
	playlist.Display.Screen = append(playlist.Display.Screen, header)
	for i := 0; i < playlist.Display.Height-2; i++ {
		rowString := fitStringToWidth("", playlist.Display.Width)
		if i < len(playlist.Playlists) {
			rowString = fitStringToWidth(playlist.Playlists[i], playlist.Display.Width)
		}
		playlist.Display.Screen = append(playlist.Display.Screen, rowString)
	}
	bottom := fillWidthWithRune('-', playlist.Display.Width)
	playlist.Display.Screen = append(playlist.Display.Screen, bottom)
}

func updateTracksDisplay(tracks *data.Tracks) {
	// TODO: Handle offet
	tracks.Display.Screen = []string{}
	header := fitStringInMiddle("Tracks", '-', tracks.Display.Width)
	tracks.Display.Screen = append(tracks.Display.Screen, header)
	for i := 0; i < tracks.Display.Height-2; i++ {
		rowString := fitStringToWidth("", tracks.Display.Width)
		if i < len(tracks.Tracks) {
			rowString = fitStringToWidth(tracks.Tracks[i], tracks.Display.Width)
		}
		tracks.Display.Screen = append(tracks.Display.Screen, rowString)
	}
	bottom := fillWidthWithRune('-', tracks.Display.Width)
	tracks.Display.Screen = append(tracks.Display.Screen, bottom)
}

// Helper Func to pad or trim string
func fitStringToWidth(str string, width int) string {
	lenStr := utf8.RuneCountInString(str)
	if lenStr <= width {
		numWhiteSpaces := width - lenStr
		return str + strings.Repeat(" ", numWhiteSpaces)
	}
	return str[:width]
}

func fitStringToWidthAndFillRune(str string, r rune, width int) string {
	lenStr := utf8.RuneCountInString(str)
	if lenStr <= width {
		numWhiteSpaces := width - lenStr
		return str + strings.Repeat(string(r), numWhiteSpaces)
	}
	return str[:width]
}

func fillWidthWithRune(r rune, width int) string {
	return strings.Repeat(string(r), width)
}

func fitStringInMiddle(str string, pad rune, width int) string {
	lenStr := utf8.RuneCountInString(str)
	if lenStr > width {
		return str[:width]
	}
	padSpace := width - lenStr
	leftPadLen := padSpace - padSpace/2
	rightPadLen := padSpace - leftPadLen
	leftPad := strings.Repeat(string(pad), leftPadLen)
	rightPad := strings.Repeat(string(pad), rightPadLen)
	return leftPad + str + rightPad
}
