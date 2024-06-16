package shell

import (
	"bytes"
	"github.com/Songmu/wrapcommander"
	"github.com/yuin/gopher-lua"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func Loader(L *lua.LState) int {
	registerLuaCommandResultType(L)

	tb := L.NewTable()
	L.SetFuncs(tb, map[string]lua.LGFunction{
		"run": run,
	})
	L.Push(tb)
	return 1
}

type CommandResult struct {
	Stdout         bytes.Buffer
	Stderr         bytes.Buffer
	CombinedOutput bytes.Buffer
	ExitStatus     int
	Err            error
}

func (r *CommandResult) Success() bool {
	return r.ExitStatus == 0
}

func (r *CommandResult) Failure() bool {
	return r.ExitStatus != 0
}

func run(L *lua.LState) int {
	command := L.CheckString(1)
	result := runCommand(command)
	L.Push(newLuaCommandResult(L, result))
	return 1
}

func runCommand(command string) *CommandResult {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var combinedOutput bytes.Buffer

	cmd.Stdout = io.MultiWriter(&stdout, &combinedOutput)
	cmd.Stderr = io.MultiWriter(&stderr, &combinedOutput)
	cmd.Stdin = os.Stdin

	var exitStatus int
	err := cmd.Run()
	if err != nil {
		exitStatus = wrapcommander.ResolveExitCode(err)
	} else {
		exitStatus = 0
	}

	return &CommandResult{
		Stdout:         stdout,
		Stderr:         stderr,
		CombinedOutput: combinedOutput,
		Err:            err,
		ExitStatus:     exitStatus,
	}
}

const luaCommandResultTypeName = "CommandResult*"

func registerLuaCommandResultType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaCommandResultTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"exit_status":     commandResultExitStatus,
		"success":         commandResultSuccess,
		"failure":         commandResultFailure,
		"stdout":          commandResultStdout,
		"stderr":          commandResultStderr,
		"combined_output": commandResultCombinedOutput,
	}))
}

func newLuaCommandResult(L *lua.LState, result *CommandResult) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = result
	L.SetMetatable(ud, L.GetTypeMetatable(luaCommandResultTypeName))
	return ud
}

func checkCommandResult(L *lua.LState) *CommandResult {
	ud := L.CheckUserData(1)
	if result, ok := ud.Value.(*CommandResult); ok {
		return result
	}
	L.ArgError(1, "CommandResult expected")
	return nil
}

func commandResultExitStatus(L *lua.LState) int {
	L.Push(lua.LNumber(checkCommandResult(L).ExitStatus))
	return 1
}

func commandResultSuccess(L *lua.LState) int {
	L.Push(lua.LBool(checkCommandResult(L).Success()))
	return 1
}

func commandResultFailure(L *lua.LState) int {
	L.Push(lua.LBool(checkCommandResult(L).Failure()))
	return 1
}

func commandResultStdout(L *lua.LState) int {
	L.Push(lua.LString(checkCommandResult(L).Stdout.String()))
	return 1
}

func commandResultStderr(L *lua.LState) int {
	L.Push(lua.LString(checkCommandResult(L).Stderr.String()))
	return 1
}

func commandResultCombinedOutput(L *lua.LState) int {
	L.Push(lua.LString(checkCommandResult(L).CombinedOutput.String()))
	return 1
}
