package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func (SpotifyPlayer) GetUserPlaylists(accessToken string) ([]SlimPlaylistData, error) {
	err := validTokenFormat(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: %w", err)
	}
	apiUrl := "https://api.spotify.com/v1/me/playlists"
	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: client: %w", err)
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, errors.New("GetUserPlaylists: Http status not successful: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: read body: %w", err)
	}

	var respStruct UserPlaylistResp

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: json: unmarshal: %w", err)
	}

	userPlaylistResp := []SlimPlaylistData{}
	for _, item := range respStruct.Items {
		newSlim := SlimPlaylistData{
			Name:         item.Name,
			DetailRefUrl: item.Href,
			TotalTracks:  item.Tracks.Total,
			TracksHref:   item.Tracks.Href,
		}
		userPlaylistResp = append(userPlaylistResp, newSlim)
	}
	return userPlaylistResp, nil
}

func (SpotifyPlayer) GetTracksFromPlaylist(hrefUrl, accessToken string, numSongs int) ([]SlimTrackInfo, error) {
	err := validTokenFormat(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: %w", err)
	}
	params := url.Values{}
	params.Add("limit", strconv.Itoa(numSongs))
	params.Add("fields", "items(track(name))")
	apiUrl := hrefUrl + "?" + params.Encode()

	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: client: %w", err)
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, errors.New("GetUserPlaylists: Http status not successful: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: read body: %w", err)
	}

	var respStruct SlimTrackResp

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("GetUserPlaylists: json: unmarshal: %w", err)
	}

	var tracks []SlimTrackInfo
	for _, item := range respStruct.Items {
		newTrack := SlimTrackInfo{
			Name: item.Track.Name,
		}
		tracks = append(tracks, newTrack)
	}

	return tracks, nil
}

func (SpotifyPlayer) GetPlaylist(hrefUrl, accessToken string) (*SlimPlaylistWithTracks, error) {
	err := validTokenFormat(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: %w", err)
	}
	err = validateUrl(hrefUrl)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: %w", err)
	}
	params := url.Values{}
	params.Add("fields", "name,tracks.items(track(name))")
	apiUrl := hrefUrl + "?" + params.Encode()

	headerStr := "Bearer " + accessToken

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: client: %w", err)
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, errors.New("GetPlaylist: Http status not successful: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: read body: %w", err)
	}

	var respStruct SlimPlaylistResp

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("GetPlaylist: json: unmarshal: %w", err)
	}
	var tracks []SlimTrackInfo
	for _, item := range respStruct.Tracks.Items {
		newTrack := SlimTrackInfo{
			Name: item.Track.Name,
		}
		tracks = append(tracks, newTrack)
	}

	slimPlaylist := SlimPlaylistWithTracks{
		PlaylistName: respStruct.Name,
		Tracks:       tracks,
	}

	return &slimPlaylist, nil
}

func validateUrl(url string) error {
	if url == "" {
		return errors.New("validateUrl: empty url")
	}
	return nil
}

type SlimPlaylistData struct {
	Name         string
	DetailRefUrl string
	TotalTracks  int
	TracksHref   string
}

type SlimTrackInfo struct {
	Name string
}

type SlimPlaylistWithTracks struct {
	PlaylistName string
	Tracks       []SlimTrackInfo
}

type UserPlaylistResp struct {
	Href     string         `json:"href"`
	Limit    int            `json:"limit"`
	Next     string         `json:"next"`
	Offset   int            `json:"offset"`
	Previous string         `json:"previous"`
	Total    int            `json:"total"`
	Items    []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"followers"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Public     bool   `json:"public"`
	SnapshotID string `json:"snapshot_id"`
	Tracks     struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type TracksFromPlaylistResp struct {
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Next     *string `json:"next"`
	Offset   int     `json:"offset"`
	Previous *string `json:"previous"`
	Total    int     `json:"total"`
	Items    []struct {
		AddedAt *string `json:"added_at"`
		AddedBy *struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"followers"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		IsLocal bool `json:"is_local"`
		Track   struct {
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
		} `json:"track"`
	} `json:"items"`
}

type SlimTrackResp struct {
	Items []struct {
		Track struct {
			Name string `json:"name"`
		} `json:"track"`
	} `json:"items"`
}

type SlimPlaylistResp struct {
	Tracks struct {
		Items []struct {
			Track struct {
				Name string `json:"name"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Name string `json:"name"`
}
