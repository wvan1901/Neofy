package input

import (
	"neofy/internal/config"
	"neofy/internal/consts"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
)

// TODO: Find a process to handle errors
func ProcessInput(d *config.AppData) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII:
		terminal.Quit(d.Term)
		break
	case 's', 'S':
		// Shuffle: FEAT:
	case 'b', 'B':
		// Previous Song
		err := spotify.SkipToPrevious(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		// TODO: Find a way to get previous song
		d.Player.CurrentSong.Name = "PrevSong"
		d.Player.CurrentSong.Artist = "PrevArtist"
	case 'p', 'P':
		// Play song
		if d.Player.IsPlaying {
			break
		}
		err := spotify.StartResumePlayback(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		d.Player.IsPlaying = true
	case 'x', 'X':
		// Pause Song
		if !d.Player.IsPlaying {
			break
		}
		err := spotify.PausePlayback(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		d.Player.IsPlaying = false
	case 'n', 'N':
		// Skip Song
		err := spotify.SkipToNext(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		// TODO: Find a way to get previous song
		d.Player.CurrentSong.Name = "NextSong"
		d.Player.CurrentSong.Artist = "NextArtist"
	case 'r', 'R':
		// Start Loop: FEAT:
	case '-':
		// Decrease Volume if enabled
		if !d.Player.SupportsVolume {
			break
		}
		newVol := d.Player.Volume - 10
		if newVol > 100 {
			newVol = 100
		} else if newVol < 0 {
			newVol = 0
		}
		err := spotify.SetPlaybackVolume(d.Spotify.UserTokens.AccessToken, newVol)
		if err != nil {
			break
		}
		d.Player.Volume = newVol
	case '+', '=':
		// Increase Volume if enabled
		if !d.Player.SupportsVolume {
			break
		}
		newVol := d.Player.Volume + 10
		if newVol > 100 {
			newVol = 100
		} else if newVol < 0 {
			newVol = 0
		}
		err := spotify.SetPlaybackVolume(d.Spotify.UserTokens.AccessToken, newVol)
		if err != nil {
			break
		}
		d.Player.Volume = newVol
	case 'f', 'F':
		// Refresh the current song
		player, err := spotify.CurrentPlayingTrack(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		d.Player.IsPlaying = player.IsPlaying
		d.Player.IsShuffled = player.IsShuffled
		d.Player.CurrentSong.Name = player.SongName
		d.Player.CurrentSong.Artist = player.Artist
		d.Player.Repeat = player.Repeat
	}
}
