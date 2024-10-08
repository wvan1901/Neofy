package mode

import (
	"neofy/internal/consts"
	"neofy/internal/data"
	"neofy/internal/terminal"
)

type Playlist struct{}

func (*Playlist) ProcessInput(d *data.AppData) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII, consts.ESC:
		d.Mode = &Player{}
		break
	case consts.CONTROL_U:
		skipBy := 10
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY-skipBy < 0 {
			break
		}
		if d.Playlist.CursorPosY-skipBy <= d.Playlist.RowOffset {
			d.Playlist.RowOffset -= skipBy
		}
		d.Playlist.CursorPosY -= skipBy
	case consts.CONTROL_D:
		skipBy := 10
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY+skipBy >= len(d.Playlist.Playlists) {
			break
		}
		if d.Playlist.CursorPosY+skipBy >= len(d.Playlist.Display.Screen)-3+d.Playlist.RowOffset {
			d.Playlist.RowOffset += skipBy
		}
		d.Playlist.CursorPosY += skipBy
	case 't', 'T':
		d.Mode = &Track{}
	case 'j', 'J':
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY+1 >= len(d.Playlist.Playlists) {
			break
		}
		if d.Playlist.CursorPosY >= len(d.Playlist.Display.Screen)-3+d.Playlist.RowOffset {
			d.Playlist.RowOffset++
		}
		d.Playlist.CursorPosY++
	case 'k', 'K':
		if d.Playlist.CursorPosY < 0 {
			break
		} else if d.Playlist.CursorPosY-1 < 0 {
			break
		}
		if d.Playlist.CursorPosY <= d.Playlist.RowOffset {
			d.Playlist.RowOffset--
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
		newTracks := []data.TrackDetail{}
		for _, track := range tracksResp {
			artists := []data.ArtistDetail{}
			for _, a := range track.Artist {
				artists = append(artists, data.ArtistDetail{Name: a.Name})
			}
			newTracks = append(newTracks, data.TrackDetail{Name: track.Name, ContextUri: track.ContextUri, DurationMs: track.DurationMs, Artists: artists})
		}
		d.Playlist.SelectedPlaylist = &curPlaylist
		d.Songs.Tracks = newTracks
		d.Songs.CursorPosY = 0
		d.Songs.SelectedTrack = nil
	}
}

func (*Playlist) ShortDisplay() rune {
	return 'U'
}
