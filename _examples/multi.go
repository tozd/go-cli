package main

import (
	"github.com/alecthomas/kong"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/zerolog"

	"gitlab.com/tozd/go/cli"
)

type Globals struct {
	zerolog.LoggingConfig `yaml:",inline"`

	Version kong.VersionFlag `help:"Show program's version and exit."                                              short:"V" yaml:"-"`
	App     cli.ConfigFlag   `help:"Load configuration from a JSON or YAML file." name:"config" placeholder:"PATH" short:"c" yaml:"-"`
}

type PlusCommand struct {
	Numbers []int `arg:"" help:"Numbers to add." name:"number" yaml:"numbers"`
}

func (c *PlusCommand) Run(globals *Globals) errors.E { //nolint:unparam
	sum := 0
	for _, n := range c.Numbers {
		sum += n
	}
	globals.Logger.Info().Msgf("%d", sum)
	return nil
}

type MinusCommand struct {
	Numbers []int `arg:"" help:"Numbers to subtract." name:"number" yaml:"numbers"`
}

func (c *MinusCommand) Run(globals *Globals) errors.E { //nolint:unparam
	// Kong makes sure there is at least one number.
	difference := c.Numbers[0]
	for _, n := range c.Numbers[1:] {
		difference -= n
	}
	globals.Logger.Info().Msgf("%d", difference)
	return nil
}

type App struct {
	Globals `yaml:"globals"`

	Plus  PlusCommand  `cmd:"" help:"Add numbers."      yaml:"plus"`
	Minus MinusCommand `cmd:"" help:"Subtract numbers." yaml:"minus"`
}

func main() {
	var app App
	cli.Run(&app, nil, func(ctx *kong.Context) errors.E {
		return errors.WithStack(ctx.Run(&app.Globals))
	})
}
