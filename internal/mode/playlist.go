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
	case 'j', 'J':
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY+1 >= len(d.Playlist.Playlists) {
			break
		}
		d.Playlist.CursorPosY++
	case 'k', 'K':
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY-1 < 0 {
			break
		}
		d.Playlist.CursorPosY--
	case 's', 'S':
		if d.Playlist.CursorPosY < 0 {
			break
		}
		if d.Playlist.SelectedPlaylist == nil {
			break
		}
		curPlaylist := d.Playlist.Playlists[d.Playlist.CursorPosY]
		tracksResp, err := d.Player.Controller.GetTracksFromPlaylist(curPlaylist.Href, d.Spotify.UserTokens.AccessToken, curPlaylist.NumSongs)
		if err != nil {
			break
		}
		newTracks := []string{}
		for _, track := range tracksResp {
			newTracks = append(newTracks, track.Name)
		}
		d.Playlist.SelectedPlaylist = &curPlaylist
		d.Songs.Tracks = newTracks
	}
}

func (*Playlist) ShortDisplay() rune {
	return 'U'
}
