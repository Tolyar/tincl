package main

import (
	. "github.com/Tolyar/tincl/internal/tincl"

	"github.com/abiosoft/ishell/v2"
)

func main() {
	var t *Telnet

	cfg := ReadConfig()

	if cfg.Interactive {
		shell := ishell.New()
		if cfg.Host != "" {
			cfg.Shell = shell
			t, _ = NewTelnet(&cfg)
		}
		shell.NotFound(TelnetInput(t))
		// Read and write history to $HOME/.ishell_history
		shell.SetHomeHistoryPath(".ishell_history")
		shell.Run()
	}
}
