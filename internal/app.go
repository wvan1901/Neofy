package internal

import (
	"fmt"
	"neofy/internal/config"
	"neofy/internal/output"
	"neofy/internal/terminal"
)

// NOTE: For our app we will have 3 panes:
// Playlist(Left), Songs(Right), Bottom(footer)
func RunApp() error {
	d := config.InitAppData()
	defer d.Term.CloseTerminal()
	fmt.Println("Access:", d.Spotify.UserTokens.AccessToken)
	//t := d.Spotify.UserTokens.AccessToken
	go d.Spotify.RefreshSchedular.Start()
	for {
		output.UpdateApp(d)
		// TODO: Implement listening to key press so we can exit the app
		terminal.ProcessKeyPress(d.Term)
	}
}
