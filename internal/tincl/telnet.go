package tincl

import (
	"bytes"
	"fmt"

	"github.com/abiosoft/ishell/v2"
	"github.com/reiver/go-telnet"
)

func TelnetInput(c *ishell.Context) {
	c.Printf("RawArgs: %s\n", c.RawArgs)
}

type Telnet struct {
	conn *telnet.Conn
}

// Read one line from telnet connection.
func (t *Telnet) ReadLine() string {
	var buf bytes.Buffer

	p := make([]byte, 1)
	buf.Reset()

	for {
		n, err := t.conn.Read(p)

		if n <= 0 && nil == err {
			continue
		} else if nil != err {
			break
		}

		if p[0] == byte('\n') {
			break
		}

		buf.Write(p)
	}

	return buf.String()
}

func NewTelnet(cfg Config) (*Telnet, error) {
	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	t := Telnet{conn: conn}

	fmt.Printf("R1: %s\n", t.ReadLine())
	conn.Write([]byte("quit\r\n"))
	fmt.Printf("R2: %s\n", t.ReadLine())
	conn.Close()

	return &t, nil
}
