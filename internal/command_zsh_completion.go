package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v2"
	"os"
	"text/template"
)

var ZshCompletionCommand = &cli.Command{
	Name:                   "zsh-completion",
	Usage:                  "Output zsh completion script to STDOUT",
	UseShortOptionHandling: true,
	Before: func(cCtx *cli.Context) error {
		// Disable debug output because it will break the zsh completion script.
		debuglogger.Get(cCtx).IsDebug = false
		return nil
	},
	Action: zshCompletionAction,
	Flags: []cli.Flag{
		zshCompletionHostsFlag,
	},
}

var zshCompletionHostsFlag = &cli.BoolFlag{Name: "hosts"}

func zshCompletionAction(cCtx *cli.Context) error {
	if zshCompletionHostsFlag.Get(cCtx) {
		return printZshCompletionHosts(cCtx)
	} else {
		return printZshCompletion(cCtx)
	}
}

func printZshCompletionHosts(cCtx *cli.Context) error {
	cfg, L, err := newConfig(cCtx)
	if err != nil {
		return err
	}
	defer L.Close()

	hosts := cfg.NewHostFilter().ExcludeHidden().GetHosts()
	for _, h := range hosts {
		_, _ = fmt.Fprintf(cCtx.App.Writer, "%s\t%s\n", h.Name, h.Description)
	}
	return nil
}

//go:embed zsh_completion.tmpl.zsh
var zshTemplateString string
var zshTmpl = template.Must(template.New("T").Parse(zshTemplateString))

func printZshCompletion(cCtx *cli.Context) error {
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

	if _, err := cCtx.App.Writer.Write(b.Bytes()); err != nil {
		return err
	}
	return nil
}
