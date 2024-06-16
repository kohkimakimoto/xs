package internal

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:                   "list",
	Aliases:                []string{"ls"},
	Usage:                  "List defined hosts",
	UseShortOptionHandling: true,
	Before: func(cCtx *cli.Context) error {
		// Disable debug output because it will break the list output.
		debuglogger.Get(cCtx).IsDebug = false
		return nil
	},
	Action: listAction,
	Flags: []cli.Flag{
		listAllFlag,
	},
}

var listAllFlag = &cli.BoolFlag{
	Name:    "all",
	Aliases: []string{"a"},
	Usage:   "List all hosts including hidden hosts",
}

func listAction(cCtx *cli.Context) error {
	cfg, L, err := newConfig(cCtx)
	if err != nil {
		return err
	}
	defer L.Close()

	f := cfg.NewHostFilter()
	if !listAllFlag.Get(cCtx) {
		f.ExcludeHidden()
	}
	hosts := f.GetHosts()

	t := newSimpleTableWriter(cCtx.App.Writer)
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
