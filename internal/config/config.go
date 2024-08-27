package config

import (
	"errors"
	"fmt"
	"neofy/internal/display"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"os"

	"github.com/joho/godotenv"
)

type AppData struct {
	Playlist Display
	Player   MusicPlayer
	Display  display.Display
	Songs    Display
	Term     terminal.AppTerm
	Spotify  spotify.Config
}

type MusicPlayer struct {
	Display        Display // What to show in cli
	IsPlaying      bool    // Is something playing
	IsShuffled     bool    // Is playlist suffled
	SupportsVolume bool    // Does Device support volume
	Volume         int     // 0-100
	CurrentSong    Song
	Repeat         string // track, context, off
}

type Display struct {
	Width  int
	Height int
	Screen []string
}

type Song struct {
	Name   string
	Artist string
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
	mp := MusicPlayer{
		Display: Display{
			Width:  newConfig.Display.Width,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
	}
	newConfig.Player = mp
	spotifyConfig, err := initSpotifyConfig()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig.Spotify = *spotifyConfig
	playerData, err := spotify.PlaybackState(spotifyConfig.UserTokens.AccessToken)
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
			Name:   "505",
			Artist: "Artic Monkeys",
		},
	}
	newConfig.Player = mp

	return &newConfig
}
