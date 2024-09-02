package internal

import (
	"fmt"
	"neofy/internal/config"
	"neofy/internal/input"
	"neofy/internal/output"
	"neofy/internal/timer"
)

// NOTE: For our app we will have 3 panes:
// Playlist(Left), Songs(Right), Bottom(footer)
func RunApp() error {
	d := config.InitAppData()
	//d := config.InitMock()
	defer d.Term.CloseTerminal()
	fmt.Println("Access:", d.Spotify.UserTokens.AccessToken)
	go d.Spotify.RefreshSchedular.Start()
	u := timer.Updater{
		Done:   make(chan bool),
		Config: d,
	}
	go u.StartWithTimer()
	for {
		output.UpdateApp(d)
		input.ProcessInput(d, &u)
	}
}
