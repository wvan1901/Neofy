package internal

import (
	"neofy/internal/config"
	"neofy/internal/data"
	"neofy/internal/output"
)

func RunApp(enableMock bool) error {
	var appData *data.AppData
	if enableMock {
		appData = config.InitMock()
	} else {
		appData = config.InitAppData()
	}

	defer appData.Term.CloseTerminal()
	go appData.Spotify.RefreshSchedular.Start()
	for {
		output.UpdateApp(appData)
		appData.Mode.ProcessInput(appData)
	}
}
