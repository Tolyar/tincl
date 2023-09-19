package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/abiosoft/ishell/v2"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const AppName = "tincl"

type Config struct {
	Interactive bool   // Run in an interactive mode.
	Host        string // Connection host.
	Port        int    // Connection port.
	Script      string // Script path.
	TLS         bool   // Use TLS connection (telnets).
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", AppName)
	pflag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nUse -i and/or -H to continue.")
	os.Exit(0)
}

func ReadConfig() Config {
	c := Config{}

	pflag.Usage = Usage
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	pflag.BoolP("interactive", "i", false, "Enable interactive mode")
	pflag.IntP("port", "P", 23, "Connection port")
	pflag.StringP("host", "H", "", "Connection host")
	pflag.StringP("script", "s", "", "Script for execution")
	pflag.BoolP("tls", "t", false, "Use TLS mode")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Can't bind flags to viper.")
	}

	// name of config file (without extension)
	viper.SetConfigName(AppName)
	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")
	// path to look for the config file in
	viper.AddConfigPath("/etc")
	// call multiple times to add many search paths
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", AppName))
	viper.AddConfigPath(".") // optionally look for config in the working directory

	// Bind ENV variables
	viper.SetEnvPrefix(AppName)
	viper.AutomaticEnv()
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			log.Fatalf("Can't read config file %v", err)
		}
	}

	c.Port = viper.GetInt("port")
	c.Host = viper.GetString("host")
	c.Interactive = viper.GetBool("interactive")
	c.Script = viper.GetString("script")

	if !c.Interactive && c.Host == "" {
		Usage()
	}

	return c
}

func TelnetInput(c *ishell.Context) {
	c.Printf("RawArgs: %s\n", c.RawArgs)
}

type Caller struct{}

func (c *Caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		oi.LongWrite(w, scanner.Bytes())
		oi.LongWrite(w, []byte("\n"))
	}
}

func OpenTelnet(cfg Config) (error, *Caller) {
	caller := Caller{}

	if cfg.TLS {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		if err := telnet.DialToAndCallTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), caller, tlsConfig); err != nil {
			return err, nil
		}
	} else {
		if err := telnet.DialToAndCall(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), caller); err != nil {
			return err, nil
		}
	}

	return nil, &caller
}

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
