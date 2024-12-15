package main

import (
	"github.com/alecthomas/kong"
	"github.com/phuslu/log"
	"github.com/toudi/kwity/cmd/commands"
	appPkg "github.com/toudi/kwity/internal/app"
)

var options struct {
	logLevel log.Level
}

var CLI struct {
	Verbosity int `type:"counter" short:"v" help:"Set verbosity level; Use -v or -vv"`

	Issue  commands.IssueCommand  `cmd:"" help:"issue invoice"`
	Render commands.RenderCommand `cmd:"" help:"render existing invoice"`
}

func main() {
	ctx := kong.Parse(&CLI)

	if CLI.Verbosity == 1 {
		options.logLevel = log.DebugLevel
	} else if CLI.Verbosity >= 2 {
		options.logLevel = log.TraceLevel
	}

	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Level:      options.logLevel,
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}

	if app, err := appPkg.NewApp(); err != nil {
		log.Fatal().Err(err).Msg("")
	} else {
		defer app.Close()
		err := ctx.Run(&commands.Context{
			App: app,
		})
		ctx.FatalIfErrorf(err)
	}
}
