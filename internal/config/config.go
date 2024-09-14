package config

import (
	"errors"
	"fmt"
	"neofy/internal/display"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"os"
	"time"
)

// TODO: Abstract Spotify & Music Player into a interface

// TODO: Add modes: Playlists, Tracks, Player
type AppData struct {
	Display  display.Display
	Mode     Mode
	Playlist Playlist
	Player   MusicPlayer
	Songs    Tracks
	Spotify  spotify.Config
	Term     terminal.AppTerm
}

// TODO: Figure out a y offset
type Playlist struct {
	SelectedPlaylist string
	Display          Display
	Playlists        []string
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

type Tracks struct {
	CurSong string
	Display Display
	Tracks  []string
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

type Mode interface {
	// TODO: How do we handle the timer?
	ProcessInput(*AppData)
	ShortDisplay() rune
}

func InitAppData() *AppData {
	newTerm := terminal.InitAppTerm()

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}

	newAppDislay := *display.InitDisplay(w, h)

	controller := spotify.SpotifyPlayer{}
	spotifyConfig, err := initSpotifyConfig()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}

	playerData, err := controller.PlaybackState(spotifyConfig.UserTokens.AccessToken)
	if err != nil {
		panic(fmt.Errorf("InitAppData: player state: %w", err))
	}

	var curSongProgress *time.Duration
	if playerData.SongProgress != nil {
		p := time.Duration(*playerData.SongProgress * 1000000)
		curSongProgress = &p
	}

	userPlaylists, err := controller.GetUserPlaylists(spotifyConfig.UserTokens.AccessToken)
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	playlists := []string{}
	for _, p := range userPlaylists {
		playlists = append(playlists, p.Name)
	}
	// TODO: Handle if the current playlist is empty
	curPlaylist, err := controller.GetPlaylist(playerData.PlaylistHref, spotifyConfig.UserTokens.AccessToken)
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	tracks := []string{}
	for _, t := range curPlaylist.Tracks {
		tracks = append(tracks, t.Name)
	}
	newPlaylist := Playlist{
		Display: Display{
			Width:  int(float64(newAppDislay.Width)*0.25) - 1,
			Height: int(float64(newAppDislay.Height)*0.9) - 1,
		},
		Playlists:        playlists,
		SelectedPlaylist: curPlaylist.PlaylistName,
	}

	mpHeight := int(float64(newAppDislay.Height) * 0.10)
	mp := MusicPlayer{
		Controller: controller,
		Display: Display{
			Width:  newAppDislay.Width - 1,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
		IsPlaying:      playerData.IsPlaying,
		SupportsVolume: playerData.SupportsVolume,
		IsShuffled:     playerData.IsShuffled,
		CurrentSong: Song{
			Name:     playerData.SongName,
			Artist:   playerData.Artist,
			Progress: curSongProgress,
			Duration: time.Duration(playerData.SongDuration * 1000000),
		},
		Repeat: playerData.Repeat,
		Volume: playerData.Volume,
	}

	newSongs := Tracks{
		CurSong: playerData.SongName,
		Display: Display{
			Width:  int(float64(newAppDislay.Width) * 0.75),
			Height: int(float64(newAppDislay.Height)*0.9) - 1,
		},
		Tracks: tracks,
	}

	newConfig := AppData{
		Display:  newAppDislay,
		Playlist: newPlaylist,
		Player:   mp,
		Songs:    newSongs,
		Spotify:  *spotifyConfig,
		Term:     newTerm,
	}

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
