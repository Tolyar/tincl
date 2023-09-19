package main

import (
	. "github.com/Tolyar/tincl/internal/tincl"

	"github.com/abiosoft/ishell/v2"
)

func main() {
	cfg := ReadConfig()

	if cfg.Host != "" {
		OpenTelnet(cfg)
	}

	if cfg.Interactive {
		shell := ishell.New()
		shell.NotFound(TelnetInput)
		// Read and write history to $HOME/.ishell_history
		shell.SetHomeHistoryPath(".ishell_history")
		shell.Run()
	}
}
