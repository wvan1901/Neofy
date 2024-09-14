package mode

import (
	"neofy/internal/consts"
	"neofy/internal/data"
	"neofy/internal/terminal"
)

type Playlist struct{}

// TODO: Implement cursor
func (*Playlist) ProcessInput(d *data.AppData) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII, consts.ESC:
		d.Mode = &Player{}
		break
	case 't', 'T':
		d.Mode = &Track{}
	}
}

func (*Playlist) ShortDisplay() rune {
	return 'U'
}
