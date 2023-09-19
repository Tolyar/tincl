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

type Caller struct{}

// func (c Caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		oi.LongWrite(w, scanner.Bytes())
// 		oi.LongWrite(w, []byte("\n"))
// 	}
// }

func (c Caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	go func(writer io.Writer, reader io.Reader) {
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
	}(os.Stdout, r)

	var buffer bytes.Buffer
	var p []byte

	crlfBuffer := [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scannerSplitFunc)

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

func scannerSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	return bufio.ScanLines(data, atEOF)
}

func OpenTelnet(cfg Config) (*Caller, error) {
	caller := Caller{}

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
