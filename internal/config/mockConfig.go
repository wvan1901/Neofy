package config

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"neofy/internal/display"
	"neofy/internal/scheduler"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"strconv"
	"time"
)

// Below is a mock config
func InitMock() *AppData {

	newTerm := terminal.InitAppTerm()

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newDisplay := *display.InitDisplay(w, h)
	mpHeight := int(float64(newDisplay.Height) * 0.10)
	progressMs := time.Millisecond * 1000 * 7
	prog := 30000
	mp := MusicPlayer{
		Display: Display{
			Width:  newDisplay.Width,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
		IsPlaying:      true,
		IsShuffled:     false,
		SupportsVolume: true,
		Volume:         77,
		Repeat:         "NONE",
		CurrentSong: Song{
			Name:     "505",
			Artist:   "Artic Monkeys",
			Duration: time.Millisecond * 1000 * 60,
			Progress: &progressMs,
		},
		Controller: &mockController{
			isPlaying:  true,
			isShuffled: false,
			volume:     77,
			repeat:     "off",
			songName:   "Init Song",
			songArtist: "Init Art",
			duration:   60000,
			progress:   &prog,
		},
	}
	newPlaylist := Playlist{
		SelectedPlaylist: "P1",
		Display: Display{
			Width:  int(float64(newDisplay.Width)*0.25) - 1,
			Height: int(float64(newDisplay.Height)*0.9) - 1,
		},
		Playlists: []string{"P1", "P2", "P3", "P4", "P5"},
	}
	newSongs := Tracks{
		Display: Display{
			Width:  int(float64(newDisplay.Width)*0.75) - 1,
			Height: int(float64(newDisplay.Height)*0.9) - 1,
		},
		Tracks: []string{"T1", "T2", "T3", "T4", "T5", "T6", "T7"},
	}
	newConfig := AppData{
		Display:  newDisplay,
		Playlist: newPlaylist,
		Player:   mp,
		Songs:    newSongs,
		Spotify:  spotify.Config{RefreshSchedular: *scheduler.CreateSchedular(time.Now(), time.Hour, nil)},
		Term:     newTerm,
	}

	return &newConfig
}

type mockController struct {
	isPlaying  bool
	isShuffled bool
	volume     int
	repeat     string
	songName   string
	songArtist string
	duration   int
	progress   *int
}

func (m *mockController) PlaybackState(string) (*spotify.SlimPlayerData, error) {
	s := spotify.SlimPlayerData{
		IsPlaying:      m.isPlaying,
		IsShuffled:     m.isShuffled,
		SupportsVolume: true,
		Volume:         m.volume,
		SongName:       m.songName,
		Artist:         m.songArtist,
		Repeat:         m.repeat,
		SongDuration:   m.duration,
		SongProgress:   m.progress,
	}
	return &s, nil
}

func (m *mockController) StartResumePlayback(string) error {
	m.isPlaying = true
	return nil
}

func (m *mockController) PausePlayback(string) error {
	m.isPlaying = false
	return nil
}

func (m *mockController) SkipToNext(string) error {
	num := rand.IntN(100)
	m.songName = "Song " + strconv.Itoa(num)
	m.songArtist = "Artist for " + strconv.Itoa(num)
	return nil
}

func (m *mockController) SkipToPrevious(string) error {
	num := rand.IntN(100) - 100
	m.songName = "Song " + strconv.Itoa(num)
	m.songArtist = "Artist for " + strconv.Itoa(num)
	return nil
}

func (m *mockController) SetPlaybackVolume(_ string, volume int) error {
	if volume > 100 {
		return nil
	} else if volume < 100 {
		return nil
	}
	m.volume = volume
	return nil
}

func (m *mockController) CurrentPlayingTrack(string) (*spotify.SlimCurrentSongData, error) {
	s := spotify.SlimCurrentSongData{
		IsPlaying:    m.isPlaying,
		IsShuffled:   m.isShuffled,
		SongName:     m.songName,
		Artist:       m.songArtist,
		Repeat:       m.repeat,
		SongDuration: m.duration,
		SongProgress: m.progress,
	}
	return &s, nil
}

func (m *mockController) RepeatMode(string, mode string) error {
	switch mode {
	case "off", "context", "track":
		m.repeat = mode
		return nil
	}
	return errors.New("RepeatMode: not valid mode")
}

func (m *mockController) ShuffleMode(_ string, b bool) error {
	m.isShuffled = b
	return nil
}

func (m *mockController) GetUserPlaylists(accessToken string) ([]spotify.SlimPlaylistData, error) {
	mockPlaylists := []spotify.SlimPlaylistData{
		{Name: "P1", DetailRefUrl: "Ref1", TotalTracks: 11, TracksHref: "t1"},
		{Name: "P2", DetailRefUrl: "Ref2", TotalTracks: 12, TracksHref: "t2"},
		{Name: "P3", DetailRefUrl: "Ref3", TotalTracks: 13, TracksHref: "t3"},
	}
	return mockPlaylists, nil
}

func (m *mockController) GetTracksFromPlaylist(string, string, int) ([]spotify.SlimTrackInfo, error) {
	mockTracks := []spotify.SlimTrackInfo{
		{Name: "Song1"}, {Name: "Song2"}, {Name: "Song3"}, {Name: "Song4"},
	}
	return mockTracks, nil
}
