package tincl

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/abiosoft/ishell/v2"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

func TelnetInput(c *ishell.Context) {
	c.Printf("RawArgs: %s\n", c.RawArgs)
}

func NewCaller(r io.ReadCloser, w io.WriteCloser) Caller {
	c := Caller{}
	c.r = r
	c.w = w

	return c
}

type Caller struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (c Caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	go func(writer io.WriteCloser, reader io.Reader) {
		var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
		p := buffer[:]

		for {
			// Read 1 byte.
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}

			oi.LongWrite(writer, p)
		}

		c.r.Close()
	}(c.w, r)

	var buffer bytes.Buffer
	var p []byte

	crlfBuffer := [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	scanner := bufio.NewScanner(c.r)
	// scanner.Split(scannerSplitFunc)

	for scanner.Scan() {
		buffer.Write(scanner.Bytes())
		buffer.Write(crlf)

		p = buffer.Bytes()

		n, err := oi.LongWrite(w, p)
		if nil != err {
			break
		}
		if expected, actual := int64(len(p)), n; expected != actual {
			err := fmt.Errorf("transmission problem: tried sending %d bytes, but actually only sent %d bytes", expected, actual)
			fmt.Fprint(os.Stderr, err.Error())

			return
		}

		buffer.Reset()
	}

	// Wait a bit to receive data from the server (that we would send to io.Stdout).
	time.Sleep(3 * time.Millisecond)
}

func OpenTelnet(cfg Config) (*Caller, error) {
	caller := NewCaller(os.Stdin, os.Stdout)

	if cfg.TLS {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		if err := telnet.DialToAndCallTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), caller, tlsConfig); err != nil {
			return nil, err
		}
	} else {
		if err := telnet.DialToAndCall(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), caller); err != nil {
			return nil, err
		}
	}

	return &caller, nil
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

		buf.Write(p)

		if p[0] == byte('\n') {
			break
		}
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
