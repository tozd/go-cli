package cli

import (
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// ConfigFlag allows you to define a config struct passed to [Kong]
// to parse a CLI argument with the path to the config file to populate
// the config struct with its contents.
//
// CLI arguments and environment variables can then override values
// populated by the config file.
//
// Config file is parsed with a YAML parser so it should be in YAML or JSON.
// Make sure the config struct supports YAML parser to populate it
// (use  "yaml" struct tags, implement [UnmarshalYAML]
// if custom parsing is needed, etc.).
//
// Example:
//
//	type App struct {
//		Config cli.ConfigFlag `short:"c" name:"config" placeholder:"PATH" help:"Load configuration from a JSON or YAML file." yaml:"-"`
//	}
//
// [Kong]: https://github.com/alecthomas/kong
// [UnmarshalYAML]: https://pkg.go.dev/gopkg.in/yaml.v3#Unmarshaler
type ConfigFlag string

func (c ConfigFlag) BeforeResolve(app *kong.Kong, ctx *kong.Context, trace *kong.Path) error {
	path := string(ctx.FlagValue(trace.Flag).(ConfigFlag)) //nolint:forcetypeassert
	file, err := os.Open(kong.ExpandPath(path))
	if err != nil {
		return errors.WithDetails(err, "path", path)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(app.Model.Target.Addr().Interface())
	if err != nil {
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			e := "error"
			if len(yamlErr.Errors) > 1 {
				e = "errors"
			}
			return errors.Errorf("yaml: unmarshal %s: %s", e, strings.Join(yamlErr.Errors, "; "))
		} else if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.WithStack(err)
	}
	return nil
}
