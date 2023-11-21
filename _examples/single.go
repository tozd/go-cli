package main

import (
	"github.com/alecthomas/kong"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/zerolog"

	"gitlab.com/tozd/go/cli"
)

const DefaultMessage = "Hello world!"

//nolint:lll
type Config struct {
	zerolog.LoggingConfig `yaml:",inline"`

	Version kong.VersionFlag `                                   help:"Show program's version and exit."                                                              short:"V" yaml:"-"`
	Config  cli.ConfigFlag   `                                   help:"Load configuration from a JSON or YAML file."   name:"config"             placeholder:"PATH"   short:"c" yaml:"-"`
	Message string           `arg:"" default:"${defaultMessage}" help:"Message to output. Default: ${defaultMessage}."               optional:"" placeholder:"STRING"           yaml:"message"`
}

func main() {
	var config Config
	cli.Run(&config, kong.Vars{
		"defaultMessage": DefaultMessage,
	}, func(ctx *kong.Context) errors.E {
		config.Logger.Info().Str("program", ctx.Model.Name).Msg(config.Message)
		return nil
	})
}
