package internal

import (
	"context"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v3"
)

func Run(args []string) error {
	return newApp().Run(context.Background(), args)
}

func newApp() *cli.Command {
	app := &cli.Command{
		Name:        "xs",
		HideVersion: true,
		Version:     Version,
		Copyright:   "Copyright (c) 2024 Kohki Makimoto",
		Metadata: map[string]any{
			"CommitHash": CommitHash,
		},
		SkipFlagParsing:               true,
		CustomRootCommandHelpTemplate: rootHelpTemplate,
	}

	app.Commands = []*cli.Command{
		SSHConfigCommand,
		ListCommand,
		ZshCompletionCommand,
		XscpFunctionCommand,
	}
	app.Before = func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		debuglogger.Bind(cmd, debuglogger.New(cmd.ErrWriter, getDebugFlag(), getNoColorFlag()))
		return ctx, nil
	}
	app.Action = func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Present() {
			first := cmd.Args().First()
			if first == "help" || first == "--help" || first == "-h" {
				return cli.ShowAppHelp(cmd)
			}
			// if args are present, run the ssh command
			return runAction(ctx, cmd)
		}

		// show help
		return cli.ShowAppHelp(cmd)
	}

	return app
}
