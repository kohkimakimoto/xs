package internal

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v3"
)

var ListCommand = &cli.Command{
	Name:                   "list",
	Aliases:                []string{"ls"},
	Usage:                  "List defined hosts",
	UseShortOptionHandling: true,
	CustomHelpTemplate:     helpTemplate,
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// Disable debug output because it will break the list output.
		debuglogger.Get(cmd).IsDebug = false
		return ctx, nil
	},
	Action: listAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "List all hosts including hidden hosts",
		},
	},
}

func listAction(ctx context.Context, cmd *cli.Command) error {
	cfg, L, err := newConfig(cmd)
	if err != nil {
		return err
	}
	defer L.Close()

	f := cfg.NewHostFilter()
	if !cmd.Bool("all") {
		f.ExcludeHidden()
	}

	hosts := f.GetHosts()

	t := newSimpleTableWriter(cmd.Writer)
	t.AppendHeader(table.Row{
		"Host",
		"Description",
		"Hidden",
	})

	for _, h := range hosts {
		t.AppendRow(table.Row{
			h.Name,
			h.Description,
			fmt.Sprintf("%t", h.Hidden),
		})
	}
	t.Render()
	return nil
}
