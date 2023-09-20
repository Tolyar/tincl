package tincl

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/abiosoft/ishell/v2"
	"github.com/reiver/go-telnet"
)

const SleepMs = 50

func TelnetInput(t *Telnet) func(c *ishell.Context) {
	return func(c *ishell.Context) {
		if t == nil {
			return
		}
		// fmt.Printf("DEBUG: \n\t%+v \n\t%+v\n", t, c)
		t.cfg.Shell.ShowPrompt(false)
		t.WriteLine(strings.Join(c.RawArgs, " "))
		// Let's chance to print something.
		time.Sleep(SleepMs * time.Millisecond)
		t.cfg.Shell.ShowPrompt(true)
	}
}

type Telnet struct {
	conn *telnet.Conn
	cfg  *Config
}

// Read one line from telnet connection.
func (t *Telnet) ReadLine() (string, error) {
	var buf bytes.Buffer

	p := make([]byte, 1)
	buf.Reset()

	for {
		n, err := t.conn.Read(p)
		if err != nil {
			return "", err
		}
		if n <= 0 {
			time.Sleep(SleepMs * time.Millisecond)

			continue
		}

		if p[0] == byte('\n') {
			break
		}

		buf.Write(p)
	}

	return buf.String(), nil
}

// Read incoming data infinite.
func (t *Telnet) readLoop() {
	for {
		t.cfg.Shell.ShowPrompt(false)

		if s, err := t.ReadLine(); err != nil {
			log.Fatalf("Connection broken by remote host %v\n", err)
			// t.cfg.Shell.ShowPrompt(true)
			break
		} else if s != "" {
			t.cfg.Shell.ShowPrompt(false)
			t.cfg.Shell.Println(s)
			t.cfg.Shell.ShowPrompt(true)
		}
		// t.cfg.Writer.Write([]byte("\n" + s + "\n"))
	}
	t.cfg.Shell.ShowPrompt(true)
}

// Read incoming data infinite.
func (t *Telnet) ReadLoop() {
	t.cfg.Shell.ShowPrompt(false)
	go t.readLoop()
	time.Sleep(SleepMs * time.Millisecond)
	t.cfg.Shell.ShowPrompt(true)
}

// Write one line to telnet.
func (t *Telnet) WriteLine(s string) (int, error) {
	return t.conn.Write([]byte(s + "\r\n"))
}

func NewTelnet(cfg *Config) (*Telnet, error) {
	cfg.Shell.Printf("* Connecting to %s:%d\n", cfg.Host, cfg.Port)
	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		cfg.Shell.Printf("ERROR: Can't connect to %s:%d - %v\n", cfg.Host, cfg.Port, err)

		return nil, err
	}

	cfg.Shell.Printf("* Successful connected to %s:%d\n", cfg.Host, cfg.Port)

	t := Telnet{conn: conn}
	t.cfg = cfg

	t.ReadLoop()

	return &t, nil
}
