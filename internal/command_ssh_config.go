package internal

import (
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v2"
)

var SSHConfigCommand = &cli.Command{
	Name:                   "ssh-config",
	Usage:                  "Output ssh_config to STDOUT",
	UseShortOptionHandling: true,
	Before: func(cCtx *cli.Context) error {
		// Disable debug output because it will break the ssh_config output.
		debuglogger.Get(cCtx).IsDebug = false
		return nil
	},
	Action: sshConfigAction,
	Flags:  []cli.Flag{},
}

func sshConfigAction(cCtx *cli.Context) error {
	cfg, L, err := newConfig(cCtx)
	if err != nil {
		return err
	}
	defer L.Close()

	sshConfigContent, err := genSSHConfig(cfg)
	if err != nil {
		return err
	}
	if _, err := cCtx.App.Writer.Write(sshConfigContent); err != nil {
		return err
	}
	return nil
}
