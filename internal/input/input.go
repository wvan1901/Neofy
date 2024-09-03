package input

import (
	"errors"
	"fmt"
	"neofy/internal/config"
	"neofy/internal/consts"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"neofy/internal/timer"
	"time"
)

// TODO: Find a process to handle errors
func ProcessInput(d *config.AppData, t *timer.Updater) {
	keyReadRune := terminal.ReadInputKey()
	switch keyReadRune {
	case consts.CONTROLCASCII:
		terminal.Quit(d.Term)
		break
	case 's', 'S':
		// Shuffle: FEAT:
		err := spotify.ShuffleMode(d.Spotify.UserTokens.AccessToken, !d.Player.IsShuffled)
		if err != nil {
			break
		}
		d.Player.IsShuffled = !d.Player.IsShuffled
	case 'b', 'B':
		// Previous Song
		err := spotify.SkipToPrevious(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		err = RefreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
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
		go t.Resume()
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
		go t.Pause()
	case 'n', 'N':
		// Skip Song
		err := spotify.SkipToNext(d.Spotify.UserTokens.AccessToken)
		if err != nil {
			break
		}
		err = RefreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
	case 'r', 'R':
		// Start Loop: FEAT:
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
		err := spotify.RepeatMode(d.Spotify.UserTokens.AccessToken, nextLoop)
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
		err := RefreshPlayer(d.Spotify.UserTokens.AccessToken, &d.Player)
		if err != nil {
			break
		}
	}
}

func RefreshPlayer(accessToken string, mp *config.MusicPlayer) error {
	if mp == nil {
		return errors.New("refreshPlayer: mp is nil")
	}
	player, err := spotify.CurrentPlayingTrack(accessToken)
	if err != nil {
		return fmt.Errorf("refreshPlayer: %w", err)
	}
	mp.IsPlaying = player.IsPlaying
	mp.IsShuffled = player.IsShuffled
	mp.CurrentSong.Name = player.SongName
	mp.CurrentSong.Artist = player.Artist
	mp.Repeat = player.Repeat
	if player.SongProgress != nil {
		p := time.Duration(*player.SongProgress * 1000000)
		mp.CurrentSong.Progress = &p
	} else {
		mp.CurrentSong.Progress = nil
	}
	mp.CurrentSong.Duration = time.Duration(player.SongDuration * 1000000)
	return nil
}
