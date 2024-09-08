package config

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"neofy/internal/display"
	"neofy/internal/scheduler"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// TODO: Abstract Spotify & Music Player into a interface

type AppData struct {
	Display  display.Display
	Playlist Display // Not Implemented Yet
	Player   MusicPlayer
	Songs    Display // Not Implemented Yet
	Spotify  spotify.Config
	Term     terminal.AppTerm
}

type MusicPlayer struct {
	Controller     spotify.Controller
	CurrentSong    Song
	Display        Display // What to show in cli
	IsPlaying      bool    // Is something playing
	IsShuffled     bool    // Is playlist suffled
	Repeat         string  // track, context, off
	SupportsVolume bool    // Does Device support volume
	Volume         int     // 0-100
}

type Display struct {
	Height int
	Screen []string
	Width  int
}

type Song struct {
	Artist   string
	Duration time.Duration
	Name     string
	Progress *time.Duration
}

func InitAppData() *AppData {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig := AppData{}

	newTerm := terminal.LinuxTerm{}
	newTerm.InitTerminal()
	newConfig.Term = &newTerm

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig.Display = *display.InitDisplay(w, h)
	pld := Display{
		Width:  int(float64(newConfig.Display.Width) * 0.25),
		Height: int(float64(newConfig.Display.Height) * 0.25),
	}
	newConfig.Playlist = pld
	sld := Display{
		Width:  int(float64(newConfig.Display.Width) * 0.75),
		Height: int(float64(newConfig.Display.Height) * 0.75),
	}
	newConfig.Songs = sld
	mpHeight := int(float64(newConfig.Display.Height) * 0.10)
	controller := spotify.SpotifyPlayer{}
	mp := MusicPlayer{
		Display: Display{
			Width:  newConfig.Display.Width,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
		Controller: controller,
	}
	newConfig.Player = mp
	spotifyConfig, err := initSpotifyConfig()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig.Spotify = *spotifyConfig
	playerData, err := controller.PlaybackState(spotifyConfig.UserTokens.AccessToken)
	if err != nil {
		panic(fmt.Errorf("InitAppData: player state: %w", err))
	}
	newConfig.Player.IsPlaying = playerData.IsPlaying
	newConfig.Player.SupportsVolume = playerData.SupportsVolume
	newConfig.Player.IsShuffled = playerData.IsShuffled
	newConfig.Player.CurrentSong = Song{
		Name:   playerData.SongName,
		Artist: playerData.Artist,
	}
	newConfig.Player.Repeat = playerData.Repeat
	newConfig.Player.Volume = playerData.Volume
	if playerData.SongProgress != nil {
		p := time.Duration(*playerData.SongProgress * 1000000)
		newConfig.Player.CurrentSong.Progress = &p
	} else {
		newConfig.Player.CurrentSong.Progress = nil
	}
	newConfig.Player.CurrentSong.Duration = time.Duration(playerData.SongDuration * 1000000)

	return &newConfig
}

func initSpotifyConfig() (*spotify.Config, error) {
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("initSpotifyConfig: ClientId is empty")
	}
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("initSpotifyConfig: ClientSecret is empty")
	}
	c := spotify.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		UserTokens:   spotify.User{},
	}
	code, err := spotify.LoginUser(clientId)
	if err != nil {
		return nil, err
	}
	accessT, refreshT, err := spotify.UserAccessAndRefreshToken(code, clientId, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("initSpotifyConfig: %w", err)
	}
	c.UserTokens.AccessToken = accessT
	c.UserTokens.RefreshToken = refreshT

	tokenScheduler := spotify.RefreshHourlyScheduler(&c.UserTokens, clientId, clientSecret)
	c.RefreshSchedular = *tokenScheduler

	return &c, nil
}

// Below is a mock config
func InitMock() *AppData {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig := AppData{}

	newTerm := terminal.LinuxTerm{}
	newTerm.InitTerminal()
	newConfig.Term = &newTerm

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig.Display = *display.InitDisplay(w, h)
	pld := Display{
		Width:  int(float64(newConfig.Display.Width) * 0.25),
		Height: int(float64(newConfig.Display.Height) * 0.9),
	}
	newConfig.Playlist = pld
	sld := Display{
		Width:  int(float64(newConfig.Display.Width) * 0.75),
		Height: int(float64(newConfig.Display.Height) * 0.9),
	}
	newConfig.Songs = sld
	mpHeight := int(float64(newConfig.Display.Height) * 0.10)
	progressMs := time.Millisecond * 1000 * 7
	prog := 30000
	mp := MusicPlayer{
		Display: Display{
			Width:  newConfig.Display.Width,
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
	newConfig.Player = mp
	tokenScheduler := scheduler.CreateSchedular(time.Now(), time.Hour, nil)
	spotifyConfig := spotify.Config{
		RefreshSchedular: *tokenScheduler,
	}
	newConfig.Spotify = spotifyConfig

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

// TODO: Finish mocks
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
