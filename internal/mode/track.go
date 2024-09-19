package mode

import (
	"neofy/internal/consts"
	"neofy/internal/data"
	"neofy/internal/terminal"
	"time"
)

type Track struct{}

func (*Track) ProcessInput(d *data.AppData) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII, consts.ESC:
		d.Mode = &Player{}
		break
	case consts.CONTROL_U:
		skipBy := 10
		if d.Songs.CursorPosY < 0 {
			break
		} else if d.Songs.CursorPosY-skipBy < 0 {
			break
		}
		if d.Songs.CursorPosY-skipBy <= d.Songs.RowOffset {
			d.Songs.RowOffset -= skipBy
		}
		d.Songs.CursorPosY -= skipBy
	case consts.CONTROL_D:
		skipBy := 10
		if d.Songs.CursorPosY < 0 {
			break
		} else if d.Songs.CursorPosY+skipBy >= len(d.Songs.Tracks) {
			break
		}
		if d.Songs.CursorPosY+skipBy >= len(d.Songs.Display.Screen)-3+d.Songs.RowOffset {
			d.Songs.RowOffset += skipBy
		}
		d.Songs.CursorPosY += skipBy
	case 'u', 'U':
		d.Mode = &Playlist{}
	case 'j', 'J':
		if d.Songs.CursorPosY < 0 {
			break
		} else if d.Songs.CursorPosY+1 >= len(d.Songs.Tracks) {
			break
		}
		if d.Songs.CursorPosY >= len(d.Songs.Display.Screen)-3+d.Songs.RowOffset {
			d.Songs.RowOffset++
		}
		d.Songs.CursorPosY++
	case 'k', 'K':
		if d.Songs.CursorPosY < 0 {
			break
		} else if d.Songs.CursorPosY-1 < 0 {
			break
		}
		if d.Songs.CursorPosY <= d.Songs.RowOffset {
			d.Songs.RowOffset--
		}
		d.Songs.CursorPosY--
	case 's', 'S':
		if d.Songs.CursorPosY < 0 {
			break
		}
		newTrack := d.Songs.Tracks[d.Songs.CursorPosY]
		err := d.Player.Controller.StartTrack(d.Playlist.SelectedPlaylist.ContextUri, d.Spotify.UserTokens.AccessToken, d.Songs.CursorPosY)
		if err != nil {
			break
		}
		artist := "???"
		if len(newTrack.Artists) > 0 {
			artist = newTrack.Artists[0].Name
		}
		zero := time.Duration(0)
		d.Songs.SelectedTrack = &newTrack
		d.Player.PlayingSong.Name = newTrack.Name
		d.Player.PlayingSong.Artist = artist
		d.Player.PlayingSong.Duration = time.Duration(newTrack.DurationMs * 1000000)
		d.Player.PlayingSong.Progress = &zero
	}
}

func (*Track) ShortDisplay() rune {
	return 'T'
}
