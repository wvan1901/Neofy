package data

import (
	"neofy/internal/display"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
	"time"
)

// TODO: Abstract Spotify & Music Player into a interface

type AppData struct {
	Display  display.Display
	Mode     Mode
	Playlist Playlist
	Player   MusicPlayer
	Songs    Tracks
	Spotify  spotify.Config
	Term     terminal.AppTerm
}

type Playlist struct {
	CursorPosY       int
	Display          Display
	Playlists        []PlaylistDetail
	RowOffset        int
	SelectedPlaylist *PlaylistDetail //Display only
}

type PlaylistDetail struct {
	Href       string
	Name       string
	NumSongs   int
	ContextUri string
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
	CursorPosY    int
	Display       Display
	RowOffset     int
	SelectedTrack *TrackDetail
	Tracks        []TrackDetail
}

type TrackDetail struct {
	Name       string
	ContextUri string
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
	ProcessInput(*AppData)
	ShortDisplay() rune
}
