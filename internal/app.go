package internal

import (
	"neofy/internal/config"
	"neofy/internal/input"
	"neofy/internal/output"
)

func RunApp() error {
	d := config.InitAppData()
	//d := config.InitMock()
	defer d.Term.CloseTerminal()
	go d.Spotify.RefreshSchedular.Start()
	for {
		output.UpdateApp(d)
		input.ProcessInput(d)
	}
}
