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

func (SpotifyPlayer) StartTrack(contextUri, accessToken string, songIndex int) error {
	err := validTokenFormat(accessToken)
	if err != nil {
		return fmt.Errorf("StartTrack: %w", err)
	}
	err = validateUrl(contextUri)
	if err != nil {
		return fmt.Errorf("StartTrack: uri: %w", err)
	}
	// NOTE: If client is currently playing then we will get a error resp
	apiUrl := "https://api.spotify.com/v1/me/player/play"
	headerStr := "Bearer " + accessToken

	reqStr := `{"context_uri": "` + contextUri + `","offset": {"position": ` + strconv.Itoa(songIndex) + `},"position_ms": 0}`
	reqBody := []byte(reqStr)
	reqBodyReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest("PUT", apiUrl, reqBodyReader)
	if err != nil {
		return fmt.Errorf("StartTrack: req: %w", err)
	}
	req.Header.Add("Authorization", headerStr)
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("StartTrack: client: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("StartTrack: read body: %w", err)
	}

	var respStruct PlayerErrorResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return fmt.Errorf("StartTrack: json: unmarshal: %w", err)
	}
	return errors.New("StartTrack: Status: " + respStruct.Error.Status + " message: " + respStruct.Error.Message)
}
