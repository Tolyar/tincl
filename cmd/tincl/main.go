package main

import (
	"bufio"
	"log"
	"os"
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
			cmd := "\r\n"
			if n, err := t.WriteLine(cmd); err != nil || n < len(cmd) {
				log.Fatalf("Can't send command to telnet. Send %d from %d bytes, with error: %v\n", n, len(cmd), err)
			}
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

	if cfg.Script != "" {
		RunScript(cfg.Script, t)
	}

	if !cfg.ReadTelnetGreeting {
		t.ReadLoop()
	}

	if cfg.Interactive {
		shell.Run()
	} else {
		scanner := bufio.NewScanner((os.Stdin))
		for scanner.Scan() {
			input := scanner.Text()
			RunCmd(t, input)
		}

	}
}
