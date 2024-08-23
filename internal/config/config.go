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
	Playlist PlayListDisplay
	Player   MusicPlayer
	Display  display.Display
	Songs    SongListDisplay
	Term     terminal.AppTerm
	Spotify  spotify.Config
}

type PlayListDisplay struct {
	Width  int
	Height int
}

type SongListDisplay struct {
	Width  int
	Height int
}

type MusicPlayer struct {
	Width  int
	Height int
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
	pld := PlayListDisplay{
		Width:  int(float64(newConfig.Display.Width) * 0.25),
		Height: int(float64(newConfig.Display.Height) * 0.25),
	}
	newConfig.Playlist = pld
	sld := SongListDisplay{
		Width:  int(float64(newConfig.Display.Width) * 0.75),
		Height: int(float64(newConfig.Display.Height) * 0.75),
	}
	newConfig.Songs = sld
	mp := MusicPlayer{
		Width:  newConfig.Display.Width,
		Height: int(float64(newConfig.Display.Height) * 0.25),
	}
	newConfig.Player = mp
	spotifyConfig, err := initSpotifyConfig()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	newConfig.Spotify = *spotifyConfig

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
