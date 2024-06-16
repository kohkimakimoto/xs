package internal

import (
	"fmt"
	"github.com/Songmu/wrapcommander"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/urfave/cli/v2"
	"github.com/yuin/gopher-lua"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func runAction(cCtx *cli.Context) error {
	logger := debuglogger.Get(cCtx)

	args := cCtx.Args().Slice()

	// extract SSH options
	var options []string
	var params []string
	for i := 0; i < len(args); i++ {
		v := args[i]
		if strings.HasPrefix(v, "-") {
			options = append(options, v)
			if v == "-B" ||
				v == "-b" ||
				v == "-c" ||
				v == "-D" ||
				v == "-E" ||
				v == "-e" ||
				v == "-F" ||
				v == "-I" ||
				v == "-i" ||
				v == "-J" ||
				v == "-L" ||
				v == "-l" ||
				v == "-m" ||
				v == "-O" ||
				v == "-o" ||
				v == "-p" ||
				v == "-Q" ||
				v == "-R" ||
				v == "-S" ||
				v == "-W" ||
				v == "-w" {
				// handle options that require values
				if i+1 < len(args) {
					options = append(options, args[i+1])
					i++
				}
			}
		} else {
			params = args[i:]
			break
		}
	}

	if len(params) == 0 {
		return fmt.Errorf("destination host is required")
	}

	tmpFile, err := os.CreateTemp("", "xs.ssh_config.*.tmp")
	if err != nil {
		return err
	}
	tmpSSHConfigFile := tmpFile.Name()
	_ = tmpFile.Close()

	defer func() {
		_ = os.Remove(tmpSSHConfigFile)
		logger.Printf("removed ssh config file: %s", tmpSSHConfigFile)
	}()

	logger.Printf("generated ssh config file: %s", tmpSSHConfigFile)

	cfg, L, err := newConfig(cCtx)
	if err != nil {
		return err
	}
	defer L.Close()

	sshConfig, err := genSSHConfig(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpSSHConfigFile, sshConfig, 0644); err != nil {
		return err
	}

	hostname := extractHostname(params[0])
	host := cfg.NewHostFilter().GetHostByName(hostname)
	if host == nil {
		return fmt.Errorf("unknown host: %s", hostname)
	}

	logger.Printf("find host: %s", host.Name)

	if len(params) == 1 {
		// If it runs without command (shell login), run hooks
		if len(host.OnBeforeConnect) > 0 {
			logger.Printf("run hooks: on_before_disconnect")
			script, err := createHookScript(L, host.OnBeforeConnect)
			if err != nil {
				return err
			}
			logger.Printf("hook script (local):")
			logger.PrintfNoPrefix("%s", script)
			if err := runHookScript(script); err != nil {
				return err
			}
		}

		if len(host.OnAfterDisconnect) > 0 {
			// register on_after_disconnect hooks
			defer func() {
				logger.Printf("run hooks: run on_after_disconnect")
				script, err := createHookScript(L, host.OnAfterDisconnect)
				if err != nil {
					_, _ = fmt.Fprintf(cCtx.App.ErrWriter, "failed to run on_after_disconnect: %v\n", err)
				}
				logger.Printf("hook script (local):")
				logger.PrintfNoPrefix("%s", script)
				if err := runHookScript(script); err != nil {
					_, _ = fmt.Fprintf(cCtx.App.ErrWriter, "failed to run on_after_disconnect: %v\n", err)
				}
			}()
		}
	}

	if len(params) == 1 && len(host.OnAfterConnect) > 0 {
		// run on_after_connect hooks
		logger.Printf("run hooks: run on_after_connect")
		script, err := createHookScript(L, host.OnAfterConnect)
		if err != nil {
			return err
		}
		// append shell login command to the end of the script
		script += "\nexec $SHELL\n"

		logger.Printf("hook script (remote):")
		logger.PrintfNoPrefix("%s", script)

		hasTOption := false
		for _, opt := range options {
			if opt == "-t" {
				hasTOption = true
				break
			}
		}
		if !hasTOption {
			// If it does not have the "-t" option, append it to the list of options.
			// The on_after_connect hook uses the ssh command with an argument to run the script.
			// This means that, by default, the ssh command does not allocate a tty.
			options = append(options, "-t")
		}

		params = append(params, script)
	}

	sshCommandArgs := []string{"-F", tmpSSHConfigFile}
	sshCommandArgs = append(sshCommandArgs, options...)
	sshCommandArgs = append(sshCommandArgs, params...)
	cmd := exec.Command("ssh", sshCommandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Printf("underlying ssh command: %v", cmd.Args)

	if err = cmd.Run(); err != nil {
		return cli.Exit(err, wrapcommander.ResolveExitCode(err))
	}
	return nil
}

func createHookScript(L *lua.LState, hooks []any) (string, error) {
	if len(hooks) == 0 {
		return "", nil
	}

	codeSlice := make([]string, 0, len(hooks))
	for _, hook := range hooks {
		code := ""
		if hookFn, ok := hook.(*lua.LFunction); ok {
			if err := L.CallByParam(lua.P{
				Fn:      hookFn,
				NRet:    1,
				Protect: true,
			}); err != nil {
				return "", err
			}

			ret := L.Get(-1) // returned value
			L.Pop(1)

			// assuming that the return value is a string
			code = lua.LVAsString(ret)
		} else if hookStr, ok := hook.(lua.LString); ok {
			code = lua.LVAsString(hookStr)
		} else {
			// never reach here if I implemented correctly.
			panic("unexpected hook type")
		}
		if code != "" {
			codeSlice = append(codeSlice, code)
		}
	}

	return strings.Join(codeSlice, "\n"), nil
}

func runHookScript(script string) error {
	if script == "" {
		return nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", script)
	} else {
		cmd = exec.Command("sh", "-c", script)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// extractHostname extracts hostname from the destination format like below:
// ssh://[user@]hostname[:port]
// [user@]hostname[:port]

func extractHostname(hostname string) string {
	hostname = strings.TrimPrefix(hostname, "ssh://")
	if strings.Contains(hostname, "@") {
		s := strings.SplitN(hostname, "@", 2)
		hostname = s[1]
	}
	if strings.Contains(hostname, ":") {
		s := strings.SplitN(hostname, ":", 2)
		hostname = s[0]
	}
	return hostname
}
