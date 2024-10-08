package config

import (
	"errors"
	"fmt"
	"neofy/internal/data"
	"neofy/internal/display"
	"neofy/internal/mode"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"os"
	"time"
)

func InitAppData() *data.AppData {
	newTerm := terminal.InitAppTerm()

	w, h, err := newTerm.GetTerminalSize()
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}

	// TODO: Figure out how to handle runes with width 2 in terminal
	newAppDislay := *display.InitDisplay(w-3, h)

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
	playlists := []data.PlaylistDetail{}
	for _, p := range userPlaylists {
		newP := data.PlaylistDetail{
			Href:       p.TracksHref,
			Name:       p.Name,
			NumSongs:   p.TotalTracks,
			ContextUri: p.ContextUri,
		}
		playlists = append(playlists, newP)
	}
	// TODO: Handle if the current playlist is empty
	curPlaylist, err := controller.GetPlaylist(playerData.PlaylistHref, spotifyConfig.UserTokens.AccessToken)
	if err != nil {
		panic(fmt.Errorf("InitAppData: %w", err))
	}
	tracks := []data.TrackDetail{}
	for _, t := range curPlaylist.Tracks {
		newT := data.TrackDetail{Name: t.Name, ContextUri: t.ContextUri}
		tracks = append(tracks, newT)
	}
	curPlaylistDetail, posY := findSelectedPlaylist(playlists, curPlaylist.PlaylistName)
	newPlaylist := data.Playlist{
		CursorPosY: posY,
		RowOffset:  0,
		Display: data.Display{
			Width:  int(float64(newAppDislay.Width)*0.25) - 1,
			Height: int(float64(newAppDislay.Height)*0.9) - 1,
		},
		Playlists:        playlists,
		SelectedPlaylist: curPlaylistDetail,
	}

	mpHeight := int(float64(newAppDislay.Height) * 0.10)
	mp := data.MusicPlayer{
		Controller: controller,
		Display: data.Display{
			Width:  newAppDislay.Width - 1,
			Height: mpHeight,
			Screen: make([]string, mpHeight),
		},
		IsPlaying:      playerData.IsPlaying,
		SupportsVolume: playerData.SupportsVolume,
		IsShuffled:     playerData.IsShuffled,
		PlayingSong: data.Song{
			Name:     playerData.SongName,
			Artist:   playerData.Artist,
			Progress: curSongProgress,
			Duration: time.Duration(playerData.SongDuration * 1000000),
		},
		Repeat: playerData.Repeat,
		Volume: playerData.Volume,
	}

	newSongs := data.Tracks{
		Display: data.Display{
			Width:  int(float64(newAppDislay.Width) * 0.75),
			Height: int(float64(newAppDislay.Height)*0.9) - 1,
		},
		SelectedTrack: &data.TrackDetail{Name: playerData.SongName},
		Tracks:        tracks,
	}

	newConfig := data.AppData{
		Display:  newAppDislay,
		Mode:     &mode.Player{},
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

func findSelectedPlaylist(list []data.PlaylistDetail, playlistName string) (*data.PlaylistDetail, int) {
	for i, p := range list {
		if p.Name == playlistName {
			return &p, i
		}
	}
	return nil, -1
}
