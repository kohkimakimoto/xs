package internal

import (
	"context"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v3"
)

var SSHConfigCommand = &cli.Command{
	Name:                   "ssh-config",
	Usage:                  "Output ssh_config to STDOUT",
	UseShortOptionHandling: true,
	CustomHelpTemplate:     helpTemplate,
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// Disable debug output because it will break the ssh_config output.
		debuglogger.Get(cmd).IsDebug = false
		return ctx, nil
	},
	Action: sshConfigAction,
	Flags:  []cli.Flag{},
}

func sshConfigAction(ctx context.Context, cmd *cli.Command) error {
	cfg, L, err := newConfig(cmd)
	if err != nil {
		return err
	}
	defer L.Close()

	sshConfigContent, err := genSSHConfig(cfg)
	if err != nil {
		return err
	}
	if _, err := cmd.Writer.Write(sshConfigContent); err != nil {
		return err
	}
	return nil
}
