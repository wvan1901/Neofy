package internal

import (
	"fmt"
	"neofy/internal/config"
	"neofy/internal/input"
	"neofy/internal/output"
)

// NOTE: For our app we will have 3 panes:
// Playlist(Left), Songs(Right), Bottom(footer)
func RunApp() error {
	//d := config.InitAppData()
	d := config.InitMock()
	defer d.Term.CloseTerminal()
	fmt.Println("Access:", d.Spotify.UserTokens.AccessToken)
	//go d.Spotify.RefreshSchedular.Start()
	for {
		output.UpdateApp(d)
		input.ProcessInput(d)
	}
}
