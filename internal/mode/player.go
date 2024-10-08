package mode

import (
	"errors"
	"fmt"
	"neofy/internal/consts"
	"neofy/internal/data"
	"neofy/internal/terminal"
	"time"
)

type Player struct{}

func (*Player) ProcessInput(d *data.AppData) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII:
		terminal.Quit(d.Term)
		break
	case 'u', 'U':
		d.Mode = &Playlist{}
	case 't', 'T':
		d.Mode = &Track{}
	case 's', 'S':
		// Shuffle:
		err := d.Player.Controller.ShuffleMode(d.Spotify.UserTokens.AccessToken, !d.Player.IsShuffled)
		if err != nil {
			break
		}
		d.Player.IsShuffled = !d.Player.IsShuffled
	case 'b', 'B':
		// Previous Song
		err := d.Player.Controller.SkipToPrevious(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		err = refreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
	case 'p', 'P':
		// Play song
		if d.Player.IsPlaying {
			break
		}
		err := d.Player.Controller.StartResumePlayback(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		d.Player.IsPlaying = true
	case 'x', 'X':
		// Pause Song
		if !d.Player.IsPlaying {
			break
		}
		err := d.Player.Controller.PausePlayback(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		d.Player.IsPlaying = false
	case 'n', 'N':
		// Skip Song
		err := d.Player.Controller.SkipToNext(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		err = refreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
	case 'r', 'R':
		// Start Loop:
		nextLoop := "off"
		switch d.Player.Repeat {
		case "off":
			nextLoop = "context"
		case "context":
			nextLoop = "track"
		case "track":
			nextLoop = "off"
		default:
			break
		}
		err := d.Player.Controller.RepeatMode(d.Spotify.UserTokens.AccessToken, nextLoop)
		if err != nil {
			break
		}
		d.Player.Repeat = nextLoop
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
		err := d.Player.Controller.SetPlaybackVolume(d.Spotify.UserTokens.AccessToken, newVol)
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
		err := d.Player.Controller.SetPlaybackVolume(d.Spotify.UserTokens.AccessToken, newVol)
		if err != nil {
			break
		}
		d.Player.Volume = newVol
	case 'f', 'F':
		// Refresh the current song
		err := refreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
	case 'w', 'W':
		panic("Wicho: Panic")
	}
}

func (*Player) ShortDisplay() rune {
	return 'P'
}

func refreshPlayer(accessToken string, mp *data.MusicPlayer) error {
	if mp == nil {
		return errors.New("refreshPlayer: mp is nil")
	}
	player, err := mp.Controller.CurrentPlayingTrack(accessToken)
	if err != nil {
		return fmt.Errorf("refreshPlayer: %w", err)
	}
	mp.IsPlaying = player.IsPlaying
	mp.IsShuffled = player.IsShuffled
	mp.PlayingSong.Name = player.SongName
	mp.PlayingSong.Artist = player.Artist
	mp.Repeat = player.Repeat
	if player.SongProgress != nil {
		p := time.Duration(*player.SongProgress * 1000000)
		mp.PlayingSong.Progress = &p
	} else {
		mp.PlayingSong.Progress = nil
	}
	mp.PlayingSong.Duration = time.Duration(player.SongDuration * 1000000)
	return nil
}
