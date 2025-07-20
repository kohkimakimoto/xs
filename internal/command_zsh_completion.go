package internal

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v3"
	"os"
	"text/template"
)

var ZshCompletionCommand = &cli.Command{
	Name:                   "zsh-completion",
	Usage:                  "Output zsh completion script to STDOUT",
	UseShortOptionHandling: true,
	CustomHelpTemplate:     helpTemplate,
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// Disable debug output because it will break the zsh completion script.
		debuglogger.Get(cmd).IsDebug = false
		return ctx, nil
	},
	Action: zshCompletionAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "hosts"},
	},
}

func zshCompletionAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("hosts") {
		return printZshCompletionHosts(ctx, cmd)
	} else {
		return printZshCompletion(ctx, cmd)
	}
}

func printZshCompletionHosts(ctx context.Context, cmd *cli.Command) error {
	cfg, L, err := newConfig(cmd)
	if err != nil {
		return err
	}
	defer L.Close()

	hosts := cfg.NewHostFilter().ExcludeHidden().GetHosts()
	for _, h := range hosts {
		_, _ = fmt.Fprintf(cmd.Writer, "%s\t%s\n", h.Name, h.Description)
	}
	return nil
}

//go:embed zsh_completion.tmpl.zsh
var zshTemplateString string
var zshTmpl = template.Must(template.New("T").Parse(zshTemplateString))

func printZshCompletion(ctx context.Context, cmd *cli.Command) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	dict := map[string]interface{}{
		"Executable": executable,
	}

	var b bytes.Buffer
	err = zshTmpl.Execute(&b, dict)
	if err != nil {
		return err
	}

	if _, err := cmd.Writer.Write(b.Bytes()); err != nil {
		return err
	}
	return nil
}
