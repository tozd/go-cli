// Package cli provides boilerplate combining [Kong] CLI argument parsing with
// [zerolog] logging.
//
// [Kong]: https://github.com/alecthomas/kong
// [zerolog]: https://gitlab.com/tozd/go/zerolog
package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/zerolog"
)

// These variables should be set during build time using "-X" ldflags.
// They are then combined into a version string and provided as [Kong] variable
// with name "version". The variable is then used by kong.VersionFlag
// to show program's version. For example, to have "-v" CLI flag show program's
// version, your Kong config struct could be:
//
//	type Config struct {
//		Version kong.VersionFlag `short:"V" help:"Show program's version and exit." json:"-" yaml:"-"`
//	}
//
// [Kong]: https://github.com/alecthomas/kong
//
//nolint:gochecknoglobals
var (
	Version        = ""
	BuildTimestamp = ""
	Revision       = ""
)

const (
	// Exit code 1 is used by Kong, 2 when program panics, 3 when program returns an error.
	errorExitCode = 3
)

type fmtError struct {
	Err error
}

func (e *fmtError) Error() string {
	return fmt.Sprintf("% -+#.1v", errors.Formatter{Error: e.Err}) //nolint:exhaustruct
}

func (e *fmtError) Unwrap() error {
	return e.Err
}

// Run runs the "run" function after [Kong] parses CLI arguments into "config" struct
// and [zerolog] logging is configured and Logger and Logger WithContext fields are set
// in "config" struct.
//
// Kong vars can override zerolog defaults and add additional variables which can then
// be interpolated in Kong struct tags in config struct. Var named "description"
// is used for program's description in usage help, if provided.
//
// Run function should always return and never call os.Exit. If it does not return
// an error, the program exits with code 0. If it returns an error, the program exits
// with code 3. The program exits with code 1 when CLI argument parsing or zerolog
// configuration fails. The program exits with code 2 on panic.
//
// Run function should not do any output to stdout by itself, but should exclusively
// use the logger. Logger then uses stdout for pretty-printed or JSON logging
// (as configured).
// Any unexpected errors go to stderr and are not in any particular format
// nor JSON (e.g., stack traces on panic).
// This combines well with [dinit].
//
// [Kong]: https://github.com/alecthomas/kong
// [zerolog]: https://gitlab.com/tozd/go/zerolog
// [dinit]: https://gitlab.com/tozd/dinit
func Run(config interface{}, vars kong.Vars, run func(*kong.Context) errors.E) {
	// Inside this function, panicking should be set to false before all regular returns from it.
	panicking := true

	parser, err := kong.New(config,
		kong.Description(vars["description"]),
		kong.UsageOnError(),
		kong.Writers(
			os.Stderr,
			os.Stderr,
		),
		kong.Vars{
			"version":                               fmt.Sprintf("version %s (build on %s, git revision %s)", Version, BuildTimestamp, Revision),
			"defaultLoggingConsoleType":             zerolog.DefaultConsoleType,
			"defaultLoggingConsoleLevel":            zerolog.DefaultConsoleLevel,
			"defaultLoggingFileLevel":               zerolog.DefaultFileLevel,
			"defaultLoggingMainLevel":               zerolog.DefaultMainLevel,
			"defaultLoggingContextLevel":            zerolog.DefaultContextLevel,
			"defaultLoggingContextConditionalLevel": zerolog.DefaultContextConditionalLevel,
			"defaultLoggingContextTriggerLevel":     zerolog.DefaultContextTriggerLevel,
		}.CloneWith(vars),
		zerolog.KongLevelTypeMapper,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: % -+#.1v", err)
		os.Exit(1)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		// We use FatalIfErrorf here because it displays usage information. But we use
		// fmtError instead of err so that we format the error and add more details
		// through its Error method, which is called inside FatalIfErrorf.
		parser.FatalIfErrorf(&fmtError{err})
	}

	// Default exist code.
	exitCode := 0
	defer func() {
		if !panicking {
			os.Exit(exitCode)
		}
	}()

	logFile, errE := zerolog.New(config)
	if logFile != nil {
		defer logFile.Close()
	}
	if errE != nil {
		parser.Fatalf("% -+#.1v", errE)
	}

	// We access main logger through global zerolog logger here, which was set in New.
	// This way we do not have to know anything about the config structure.
	logger := log.Logger

	errE = run(ctx)
	if errE != nil {
		logger.Error().Err(errE).Send()
		exitCode = errorExitCode
	}

	panicking = false
}
