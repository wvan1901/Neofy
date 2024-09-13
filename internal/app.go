package internal

import (
	"neofy/internal/config"
	"neofy/internal/input"
	"neofy/internal/output"
	"neofy/internal/timer"
)

func RunApp() error {
	d := config.InitAppData()
	//d := config.InitMock()
	defer d.Term.CloseTerminal()
	go d.Spotify.RefreshSchedular.Start()
	songDataUpdater := timer.Updater{
		Done:   make(chan bool),
		Config: d,
	}
	go songDataUpdater.StartWithTimer()
	for {
		output.UpdateApp(d)
		input.ProcessInput(d, &songDataUpdater)
	}
}
