package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

type Controller interface {
	PlaybackState(string) (*SlimPlayerData, error)
	StartResumePlayback(string) error
	PausePlayback(string) error
	SkipToNext(string) error
	SkipToPrevious(string) error
	SetPlaybackVolume(string, int) error
	CurrentPlayingTrack(string) (*SlimCurrentSongData, error)
	RepeatMode(string, string) error
	ShuffleMode(string, bool) error
	GetUserPlaylists(string) ([]SlimPlaylistData, error)
	GetTracksFromPlaylist(string, string, int) ([]SlimTrackInfo, error)
}

type SpotifyPlayer struct{}

// NOTE: If we get a 401 it most likely means we didnt request permission
// when asking user access token, add scope to func: AuthorizeUserUrl()

// NOTE: Player api endppoints are only avaliable for spotify premium members

func (SpotifyPlayer) PlaybackState(accessToken string) (*SlimPlayerData, error) {
	err := validTokenFormat(accessToken)
	if err != nil {
		return nil, fmt.Errorf("PlaybackState: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/player"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("PlaybackState: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PlaybackState: client: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("PlaybackState: resp code isn't 200 got: " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("PlaybackState: read body: %w", err)
	}
	var respStruct playbackStateResponse

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("PlaybackState: json: unmarshal: %w", err)
	}

	slimResp := SlimPlayerData{
		IsPlaying:      respStruct.IsPlaying,
		IsShuffled:     respStruct.ShuffleState,
		SupportsVolume: respStruct.Device.SupportsVolume,
		Volume:         *respStruct.Device.VolumePercent,
		SongName:       respStruct.Item.Name,
		Artist:         respStruct.Item.Artists[0].Name,
		Repeat:         respStruct.RepeatState,
		SongProgress:   respStruct.ProgressMs,
		SongDuration:   respStruct.Item.DurationMs,
	}

	return &slimResp, nil
}

func (SpotifyPlayer) StartResumePlayback(accessToken string) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("StartResumePlayback: %w", err)
	}
	// NOTE: If client is currently playing then we will get a error resp
	apiUrl := "https://api.spotify.com/v1/me/player/play"
	headerStr := "Bearer " + accessToken

	reqBody := []byte(`{"position_ms": 0}`)
	reqBodyReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest("PUT", apiUrl, reqBodyReader)
	if err != nil {
		return fmt.Errorf("StartResumePlayback: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("StartResumePlayback: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("StartResumePlayback: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("StartResumePlayback: json: unmarshal: %w", err)
	}
	return errors.New("StartResumePlayback: Status: " + respStruct.Error.Status + " message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) PausePlayback(accessToken string) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("PausePlayback: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/player/pause"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("PausePlayback: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("PausePlayback: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("PausePlayback: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("PausePlayback: json: unmarshal: %w", err)
	}
	return errors.New("PausePlayback: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) SkipToNext(accessToken string) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("SkipToNext: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/player/next"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("SkipToNext: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("SkipToNext: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SkipToNext: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("SkipToNext: json: unmarshal: %w", err)
	}
	return errors.New("SkipToNext: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) SkipToPrevious(accessToken string) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("SkipToPrevious: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/player/previous"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("SkipToPrevious: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("SkipToPrevious: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SkipToPrevious: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("SkipToPrevious: json: unmarshal: %w", err)
	}
	return errors.New("SkipToPrevious: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) SetPlaybackVolume(accessToken string, volume int) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("SetPlaybackVolume: %w", err)
	}
	if volume < 0 || volume > 100 {
		return errors.New("SetPlaybackVolume: Volume must be between 0-100")
	}
	apiUrl := "https://api.spotify.com/v1/me/player/volume?volume_percent=" + strconv.Itoa(volume)
	headerStr := "Bearer " + accessToken
	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("SetPlaybackVolume: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("SetPlaybackVolume: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SetPlaybackVolume: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("SetPlaybackVolume: json: unmarshal: %w", err)
	}
	return errors.New("SetPlaybackVolume: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) CurrentPlayingTrack(accessToken string) (*SlimCurrentSongData, error) {
	err := validTokenFormat(accessToken)
	if err != nil {
		return nil, fmt.Errorf("CurrentPlayingTrack: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/player/currently-playing"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("CurrentPlayingTrack: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CurrentPlayingTrack: client: %w", err)
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, errors.New("CurrentPlayingTrack: Http status not successful: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CurrentPlayingTrack: read body: %w", err)
	}

	var respStruct currentTrackResponse

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("CurrentPlayingTrack: json: unmarshal: %w", err)
	}

	slimResp := SlimCurrentSongData{
		IsPlaying:    respStruct.IsPlaying,
		IsShuffled:   respStruct.ShuffleState,
		SongName:     respStruct.Item.Name,
		Artist:       respStruct.Item.Artists[0].Name,
		Repeat:       respStruct.RepeatState,
		SongProgress: respStruct.ProgressMs,
		SongDuration: respStruct.Item.DurationMs,
	}
	return &slimResp, nil
}

func (SpotifyPlayer) RepeatMode(accessToken string, state string) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("SetRepeatMode: %w", err)
	}
	options := []string{"off", "context", "track"}
	if !slices.Contains(options, state) {
		return errors.New("SetRepeatMode: not a valid state")
	}

	apiUrl := "https://api.spotify.com/v1/me/player/repeat?state=" + state
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("SetRepeatMode: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("SetRepeatMode: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SetRepeatMode: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("SetRepeatMode: json: unmarshal: %w", err)
	}
	return errors.New("SetRepeatMode: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func (SpotifyPlayer) ShuffleMode(accessToken string, isShuffled bool) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("ShuffleMode: %w", err)
	}
	shuffled := "false"
	if isShuffled {
		shuffled = "true"
	}

	apiUrl := "https://api.spotify.com/v1/me/player/shuffle?state=" + shuffled
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("ShuffleMode: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("ShuffleMode: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ShuffleMode: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("ShuffleMode: json: unmarshal: %w", err)
	}
	return errors.New("ShuffleMode: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func validTokenFormat(token string) error {
	if token == "" {
		return errors.New("validTokenFormat: empty access token")
	}
	return nil
}

// Structs:
type SlimPlayerData struct {
	IsPlaying      bool
	IsShuffled     bool
	SupportsVolume bool
	Volume         int
	SongName       string
	Artist         string
	Repeat         string
	SongDuration   int
	SongProgress   *int
}

type SlimCurrentSongData struct {
	IsPlaying    bool
	IsShuffled   bool
	SongName     string
	Artist       string
	Repeat       string
	SongDuration int
	SongProgress *int
}

// playbackStateResponse
type playbackStateResponse struct {
	Device               Device  `json:"device"`
	RepeatState          string  `json:"repeat_state"`
	ShuffleState         bool    `json:"shuffle_state"`
	Context              Context `json:"context"`
	Timestamp            int     `json:"timestamp"`
	ProgressMs           *int    `json:"progress_ms"`
	IsPlaying            bool    `json:"is_playing"`
	Item                 Item    `json:"item"`
	CurrentlyPlayingType string  `json:"currently_playing_type"`
	Actions              Action  `json:"actions"`
}

type currentTrackResponse struct {
	RepeatState          string            `json:"repeat_state"`
	ShuffleState         bool              `json:"shuffle_state"`
	Context              Context           `json:"context"`
	Timestamp            int               `json:"timestamp"`
	ProgressMs           *int              `json:"progress_ms"`
	IsPlaying            bool              `json:"is_playing"`
	Item                 Item              `json:"item"`
	CurrentlyPlayingType string            `json:"currently_playing_type"`
	Actions              CurrentSongAction `json:"actions"`
}
type Device struct {
	ID               *string `json:"id"`
	IsActive         bool    `json:"is_active"`
	IsPrivateSession bool    `json:"is_private_session"`
	IsRestricted     bool    `json:"is_restricted"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	VolumePercent    *int    `json:"volume_percent"`
	SupportsVolume   bool    `json:"supports_volume"`
}

type Context struct {
	Type         string `json:"type"`
	Href         string `json:"href"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	URI string `json:"uri"`
}

type Item struct {
	Album struct {
		AlbumType        string   `json:"album_type"`
		TotalTracks      int      `json:"total_tracks"`
		AvailableMarkets []string `json:"available_markets"`
		ExternalUrls     struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height *int   `json:"height"`
			Width  *int   `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		Restrictions         struct {
			Reason string `json:"reason"`
		} `json:"restrictions"`
		Type    string `json:"type"`
		URI     string `json:"uri"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
		Ean  string `json:"ean"`
		Upc  string `json:"upc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href       string `json:"href"`
	ID         string `json:"id"`
	IsPlayable bool   `json:"is_playable"`
	LinkedFrom struct {
	} `json:"linked_from"`
	Restrictions struct {
		Reason string `json:"reason"`
	} `json:"restrictions"`
	Name        string  `json:"name"`
	Popularity  int     `json:"popularity"`
	PreviewURL  *string `json:"preview_url"`
	TrackNumber int     `json:"track_number"`
	Type        string  `json:"type"`
	URI         string  `json:"uri"`
	IsLocal     bool    `json:"is_local"`
}

type Action struct {
	InterruptingPlayback  *bool `json:"interrupting_playback"`
	Pausing               *bool `json:"pausing"`
	Resuming              *bool `json:"resuming"`
	Seeking               *bool `json:"seeking"`
	SkippingNext          *bool `json:"skipping_next"`
	SkippingPrev          *bool `json:"skipping_prev"`
	TogglingRepeatContext *bool `json:"toggling_repeat_context"`
	TogglingShuffle       *bool `json:"toggling_shuffle"`
	TogglingRepeatTrack   *bool `json:"toggling_repeat_track"`
	TransferringPlayback  *bool `json:"transferring_playback"`
}

type CurrentSongAction struct {
	Disallows struct {
		Resuming bool `json:"resuming"`
	} `json:"disallows"`
}

type PlayerErrorResponse struct {
	Error struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}
