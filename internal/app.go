package internal

import (
	"fmt"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v2"
)

func Run(args []string) error {
	return newApp().Run(args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.HideHelp = true
	app.HideVersion = true
	app.Name = "XS"
	app.Version = Version
	app.Copyright = "Copyright (c) 2024 Kohki Makimoto"
	app.Metadata = map[string]interface{}{
		"CommitHash": CommitHash,
	}
	app.SkipFlagParsing = true
	app.Commands = []*cli.Command{
		SSHConfigCommand,
		ListCommand,
		ZshCompletionCommand,
		XscpFunctionCommand,
	}
	app.Before = func(cCtx *cli.Context) error {
		debuglogger.Bind(cCtx.App, debuglogger.New(cCtx.App.ErrWriter, getDebugFlag(), getNoColorFlag()))
		return nil
	}
	app.Action = func(cCtx *cli.Context) error {
		if cCtx.Args().Present() {
			first := cCtx.Args().First()
			if first == "help" || first == "--help" || first == "-h" {
				return cli.ShowAppHelp(cCtx)
			}
			if first == "-V" {
				_, _ = fmt.Fprintf(cCtx.App.Writer, "%s %s\n", cCtx.App.Name, cCtx.App.Version)
				return nil
			}

			// if args are present, run the ssh command
			return runAction(cCtx)
		}

		// show help
		return cli.ShowAppHelp(cCtx)
	}

	return app
}

func init() {
	cli.AppHelpTemplate = `Usage: xs [options] builtin_command|destination [command [args ...]]

XS is a SSH command wrapper that enhances your SSH operations.

Options:
   -h, --help     Show this help message and exit
   You can also use ssh command options. Check 'man ssh' for more information.

Builtin commands:{{template "visibleCommandCategoryTemplate" .}}

Destination Hosts:
   You can define destination hosts in the configuration file.

Environment variables:
   XS_CONFIG_FILE  Path to the configuration file. Default is ~/.xs/config.lua
   XS_DEBUG        If set to "true", XS will output debug information.
   XS_NO_COLOR     If set to "true", XS will not output color codes in debug information.

Version: {{ .Version }}
Commit: {{ .Metadata.CommitHash }}
{{template "copyrightTemplate" .}}
`
	cli.CommandHelpTemplate = `Usage: {{template "usageTemplate" .}}

{{template "helpNameTemplate" .}}

Options:{{template "visibleFlagTemplate" .}}
`
}
