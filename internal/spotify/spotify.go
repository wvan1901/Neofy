package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"neofy/internal/scheduler"
	"neofy/internal/terminal"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	REDIRECT_URI = "http://localhost:8090/callback"
)

// TODO: Rewrite config & make it into interface to mock api calls

type Config struct {
	ClientId         string
	ClientSecret     string
	UserTokens       User
	RefreshSchedular scheduler.Schedular
}

type User struct {
	AccessToken  string
	RefreshToken string
}

func AccessToken(clientId, clientSecret string) (string, error) {
	apiUrl := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_secret", clientSecret)
	data.Add("client_id", clientId)
	postData := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiUrl, postData)
	if err != nil {
		return "", fmt.Errorf("AccessToken: req: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("AccessToken: client: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("AccessToken: read body: %w", err)
	}
	respStruct := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "", fmt.Errorf("AccessToken: json: unmarshal: %w", err)
	}

	return respStruct.AccessToken, nil
}

// Return a url that the user will use to auth for spotify
func AuthorizeUserUrl(clientId string) (string, error) {
	apiUrl := "https://accounts.spotify.com/authorize"
	redirectUri := REDIRECT_URI
	data := url.Values{}
	data.Add("client_id", clientId)
	data.Add("response_type", "code")
	data.Add("redirect_uri", redirectUri)
	data.Add("scope", "user-modify-playback-state user-read-playback-state playlist-read-private")
	reqUrl := apiUrl + "?" + data.Encode()

	return reqUrl, nil
}

// NOTE: FEATURE: Make this better so we can avoid chanels blocking & panics
func LoginUser(clientId string) (string, error) {
	// Get url for user to auth
	userLoginUrl, err := AuthorizeUserUrl(clientId)
	if err != nil {
		return "", fmt.Errorf("LoginUser: %w", err)
	}
	// Start Server to listen for callback then kill it
	codeChan := make(chan string)
	loginSrv := CreateLoginServer(codeChan)
	go loginSrv.RunServer()
	// Auto open link for user
	err = terminal.Openbrowser(userLoginUrl)
	// Use callback to extract the code
	code := <-codeChan
	time.Sleep(2 * time.Second)
	close(codeChan)
	// User is logged in by having a code
	return code, nil
}

func UserAccessAndRefreshToken(code, clientId, clientSecret string) (string, string, error) {
	apiUrl := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", REDIRECT_URI)
	postData := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiUrl, postData)
	if err != nil {
		return "", "", fmt.Errorf("UserAccessAndRefreshToken: req: %w", err)
	}
	authString := clientId + ":" + clientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Add("Authorization", encodedAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("UserAccessAndRefreshToken: client: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("UserAccessAndRefreshToken: read body: %w", err)
	}
	respStruct := struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "", "", fmt.Errorf("UserAccessAndRefreshToken: json: unmarshal: %w", err)
	}

	return respStruct.AccessToken, respStruct.RefreshToken, nil
}

func RefreshUserTokens(refreshToken, clientId, clientSecret string) (string, string, error) {
	apiUrl := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	postData := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiUrl, postData)
	if err != nil {
		return "", "", fmt.Errorf("RefreshUserTokens: req: %w", err)
	}
	authString := clientId + ":" + clientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Add("Authorization", encodedAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("RefreshUserTokens: client: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("RefreshUserTokens: read body: %w", err)
	}
	respStruct := struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "", "", fmt.Errorf("RefreshUserTokens: json: unmarshal: %w", err)
	}

	return respStruct.AccessToken, respStruct.RefreshToken, nil
}

// Hourly Scheduler to request new tokens & save it
type refreshTokenJob struct {
	userTokens   *User
	clientId     string
	clientSecret string
}

func (s *refreshTokenJob) Execute() {
	// Call refresh token API
	newAccess, newRefresh, err := RefreshUserTokens(s.userTokens.RefreshToken, s.clientId, s.clientSecret)
	if err != nil {
		// Handle the error
		panic(fmt.Errorf("Refresh hourly scheduler: Execute: %w", err))
	}
	// Update value
	s.userTokens.AccessToken = newAccess
	// NOTE: When a refresh token is not returned, continue using the existing token.
	if newRefresh == "" {
		newRefresh = s.userTokens.RefreshToken
	}
	s.userTokens.RefreshToken = newRefresh
}

func RefreshHourlyScheduler(u *User, clientId, clientSecret string) *scheduler.Schedular {
	startTime := time.Now().Add(time.Minute * 57)
	delay := time.Minute * 57
	j := refreshTokenJob{
		userTokens:   u,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
	var js []scheduler.Job
	js = append(js, &j)
	s := scheduler.CreateSchedular(startTime, delay, js)
	return s
}
