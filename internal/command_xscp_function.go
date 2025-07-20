package internal

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v3"
	"os"
	"text/template"
)

var XscpFunctionCommand = &cli.Command{
	Name:                   "xscp-function",
	Usage:                  "Output xscp function code to STDOUT",
	UseShortOptionHandling: true,
	CustomHelpTemplate:     helpTemplate,
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// Disable debug output because it will break the script.
		debuglogger.Get(cmd).IsDebug = false
		return ctx, nil
	},
	Action: xscpFunctionAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			DefaultText: "xscp",
			Usage:       "function name",
		},
	},
}

//go:embed xscp_function.tmpl.sh
var xscpTemplateString string
var xscpTmpl = template.Must(template.New("T").Parse(xscpTemplateString))

func xscpFunctionAction(ctx context.Context, cmd *cli.Command) error {
	name := cmd.String("name")
	if name == "" {
		name = "xscp"
	}

	executable, err := os.Executable()
	if err != nil {
		return err
	}
	dict := map[string]interface{}{
		"Executable": executable,
		"Name":       name,
	}

	var b bytes.Buffer
	if err := xscpTmpl.Execute(&b, dict); err != nil {
		return err
	}

	if _, err := cmd.Writer.Write(b.Bytes()); err != nil {
		return err
	}
	return nil
}
