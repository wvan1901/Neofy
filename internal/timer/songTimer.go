package timer

import (
	"fmt"
	"neofy/internal/config"
	"neofy/internal/output"
	"neofy/internal/spotify"
	"time"
)

// Requirements: Start, Pause, Other
type Updater struct {
	Done   chan bool
	Config *config.AppData
	tic    *time.Ticker
}

func (s *Updater) update() error {
	err := refreshPlayerData(s.Config)
	if err != nil {
		return err
	}

	output.UpdateApp(s.Config)
	return nil
}

func (s *Updater) Stop() {
	s.Done <- true
}

func (s *Updater) Pause() {
	s.tic.Reset(time.Hour)
}

func (s *Updater) Resume() {
	err := refreshPlayerData(s.Config)
	if err != nil {
		return
	}
	progress := time.Second * 0
	duration := time.Hour
	if s.Config.Player.CurrentSong.Progress != nil {
		progress = *s.Config.Player.CurrentSong.Progress
		duration = s.Config.Player.CurrentSong.Duration
	}
	s.tic.Reset(duration - progress + time.Millisecond*100)
}

func (s *Updater) StartWithTimer() {
	progress := time.Second * 0
	duration := time.Hour
	if s.Config.Player.CurrentSong.Progress != nil {
		progress = *s.Config.Player.CurrentSong.Progress
		duration = s.Config.Player.CurrentSong.Duration
	}
	s.tic = time.NewTicker(duration - progress + time.Millisecond*100)
	defer s.tic.Stop()

	for {
		select {
		case <-s.Done:
			return
		case <-s.tic.C:
			err := s.update()
			if err != nil {
				continue
			}
			newProg := time.Second * 0
			newDur := time.Hour
			if s.Config.Player.CurrentSong.Progress != nil {
				newProg = *s.Config.Player.CurrentSong.Progress
				newDur = s.Config.Player.CurrentSong.Duration
			}
			s.tic.Reset(newDur - newProg + time.Millisecond*100)
		}
	}
}

func (s *Updater) Kill() {
	close(s.Done)
}

func refreshPlayerData(d *config.AppData) error {
	player, err := spotify.CurrentPlayingTrack(d.Spotify.UserTokens.AccessToken)
	if err != nil {
		return fmt.Errorf("refreshPlayer: %w", err)
	}
	d.Player.IsPlaying = player.IsPlaying
	d.Player.IsShuffled = player.IsShuffled
	d.Player.CurrentSong.Name = player.SongName
	d.Player.CurrentSong.Artist = player.Artist
	d.Player.Repeat = player.Repeat
	if player.SongProgress != nil {
		p := time.Duration(*player.SongProgress * 1000000)
		d.Player.CurrentSong.Progress = &p
	} else {
		d.Player.CurrentSong.Progress = nil
	}
	d.Player.CurrentSong.Duration = time.Duration(player.SongDuration * 1000000)
	return nil
}
