package mode

import (
	"neofy/internal/consts"
	"neofy/internal/data"
	"neofy/internal/terminal"
)

type Track struct{}

// TODO: Implement cursor
func (*Track) ProcessInput(d *data.AppData) {
	// TODO: Find a way to remove timer, play & pause use it
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII, consts.ESC:
		d.Mode = &Player{}
		break
	case 'u', 'U':
		d.Mode = &Playlist{}
	}
}

func (*Track) ShortDisplay() rune {
	return 'T'
}
