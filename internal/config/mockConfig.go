package config

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"neofy/internal/data"
	"neofy/internal/display"
	"neofy/internal/mode"
	"neofy/internal/scheduler"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"strconv"
	"time"
)

// Below is a mock config
func InitMock() *data.AppData {

	newTerm := terminal.InitAppTerm()

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newDisplay := *display.InitDisplay(w, h)
	mpHeight := int(float64(newDisplay.Height) * 0.10)
	progressMs := time.Millisecond * 1000 * 7
	prog := 30000
	mp := data.MusicPlayer{
		Display: data.Display{
			Width:  newDisplay.Width - 2,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
		IsPlaying:      true,
		IsShuffled:     false,
		SupportsVolume: true,
		Volume:         77,
		Repeat:         "NONE",
		CurrentSong: data.Song{
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
	playlists := []data.PlaylistDetail{{Name: "P1"}, {Name: "P2"}, {Name: "P3"}, {Name: "P4"}, {Name: "P5"}}
	posY := 0
	curPlaylist := playlists[posY]
	newPlaylist := data.Playlist{
		CursorPosY:       posY,
		RowOffset:        0,
		SelectedPlaylist: &curPlaylist,
		Display: data.Display{
			Width:  int(float64(newDisplay.Width)*0.25) - 1,
			Height: int(float64(newDisplay.Height)*0.9) - 1,
		},
		Playlists: playlists,
	}
	newSongs := data.Tracks{
		Display: data.Display{
			Width:  int(float64(newDisplay.Width)*0.75) - 1,
			Height: int(float64(newDisplay.Height)*0.9) - 1,
		},
		Tracks:    []data.TrackDetail{{Name: "T1"}, {Name: "T2"}, {Name: "T3"}},
		RowOffset: 0,
	}
	newConfig := data.AppData{
		Display:  newDisplay,
		Mode:     &mode.Player{},
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
	randLen := rand.IntN(10) + 1
	mocks := []spotify.SlimTrackInfo{}
	for i := 1; i <= randLen; i++ {
		mocks = append(mocks, spotify.SlimTrackInfo{Name: "Song" + strconv.Itoa(i)})
	}
	return mocks, nil
}
