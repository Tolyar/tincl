package tincl

import (
	"fmt"
	"log"
	"os"

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
