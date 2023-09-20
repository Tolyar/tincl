package main

import (
	"log"
	"time"

	. "github.com/Tolyar/tincl/internal/tincl"

	"github.com/abiosoft/ishell/v2"
)

func main() {
	var t *Telnet
	var err error

	cfg := ReadConfig()
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "crlf",
		Help: "Send crlf",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			t.WriteLine("\r\n")
			time.Sleep(50 * time.Millisecond)
			c.ShowPrompt(true)
		},
	})

	if cfg.Host != "" {
		cfg.Shell = shell
		t, err = NewTelnet(&cfg)
		if err != nil {
			log.Fatalf("* Can't connect to remote server: %v\n", err)
		}
	}
	shell.NotFound(TelnetInput(t))
	// Read and write history to $HOME/.ishell_history
	shell.SetHomeHistoryPath(".ishell_history")

	if cfg.Interactive {
		shell.Run()
	}
}
