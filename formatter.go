package cli

import (
	"strings"

	"github.com/alecthomas/kong"
)

// DefaultValueFormatter is a HelpValueFormatter which appends automatically
// possible enum values, default values and environment variables to the help.
func DefaultValueFormatter(value *kong.Value) string {
	help := value.Help
	if value.Enum != "" {
		help += " Possible: " + value.Enum + "."
	}
	if value.Default != "" {
		help += " Default: " + value.Default + "."
	}
	if len(value.Tag.Envs) > 1 {
		help += " Environment variables: " + strings.Join(value.Tag.Envs, ",") + "."
	} else if len(value.Tag.Envs) == 1 {
		help += " Environment variable: " + value.Tag.Envs[0] + "."
	}
	return help
}
