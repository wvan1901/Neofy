package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Handle non 2** http responses

// TODO: Fix error string messages

// NOTE: If we get a 401 it most likely means we didnt request permission
// when asking user access token, add scope to func: AuthorizeUserUrl()

// NOTE: When we send a api request we should block/freeze app until we get a response

// NOTE: Only player api endppoint is only avaliable for spotify premium members

func PlaybackState(accessToken string) error {
	// TODO: Figure out what information is useful for app
	apiUrl := "https://api.spotify.com/v1/me/player"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("PlaybackState: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("PlaybackState: client: %w", err)
	}
	if resp.StatusCode != 200 {
		return errors.New("PlaybackState: resp code isn't 200 got: " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("PlaybackState: read body: %w", err)
	}
	var respStruct playbackStateResponse

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("PlaybackState: json: unmarshal: %w", err)
	}
	fmt.Println("Wicho: PlaybackState: isPlaying:", respStruct.IsPlaying)

	return nil
}

func StartResumePlayback(accessToken string) error {
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
	return errors.New("StartResumePlayback: Status: " + respStruct.Error.Status + "message: " + respStruct.Error.Message)
}

func PausePlayback(accessToken string) error {
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

func SkipToNext(accessToken string) error {
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

func SkipToPrevious(accessToken string) error {
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

func SetPlaybackVolume(accessToken string, volume int) error {
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

// Structs:
type playbackStateResponse struct {
	Device struct {
		ID               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
		SupportsVolume   bool   `json:"supports_volume"`
	} `json:"device"`
	RepeatState  string `json:"repeat_state"`
	ShuffleState bool   `json:"shuffle_state"`
	Context      struct {
		Type         string `json:"type"`
		Href         string `json:"href"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		URI string `json:"uri"`
	} `json:"context"`
	Timestamp  int  `json:"timestamp"`
	ProgressMs int  `json:"progress_ms"`
	IsPlaying  bool `json:"is_playing"`
	Item       struct {
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
				Height int    `json:"height"`
				Width  int    `json:"width"`
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
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
		IsLocal     bool   `json:"is_local"`
	} `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions              struct {
		InterruptingPlayback  bool `json:"interrupting_playback"`
		Pausing               bool `json:"pausing"`
		Resuming              bool `json:"resuming"`
		Seeking               bool `json:"seeking"`
		SkippingNext          bool `json:"skipping_next"`
		SkippingPrev          bool `json:"skipping_prev"`
		TogglingRepeatContext bool `json:"toggling_repeat_context"`
		TogglingShuffle       bool `json:"toggling_shuffle"`
		TogglingRepeatTrack   bool `json:"toggling_repeat_track"`
		TransferringPlayback  bool `json:"transferring_playback"`
	} `json:"actions"`
}

type PlayerErrorResponse struct {
	Error struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}
