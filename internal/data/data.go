package data

import (
	"neofy/internal/display"
	"neofy/internal/spotify"
	"neofy/internal/terminal"
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
	ProcessInput(*AppData) // Needs Terminal, Player, Spotify
	ShortDisplay() rune
}
