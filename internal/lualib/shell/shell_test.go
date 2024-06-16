package shell

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()
		L.PreloadModule("shell", Loader)

		code := `
local shell = require("shell")

local result = shell.run("echo hello")
assert(result:exit_status() == 0)
assert(result:success() == true)
assert(result:failure() == false)
assert(result:stdout() == "hello\n")
assert(result:stderr() == "")
assert(result:combined_output() == "hello\n")
`
		err := L.DoString(code)
		assert.NoError(t, err)
	})

	t.Run("run with error", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()
		L.PreloadModule("shell", Loader)

		code := `
local shell = require("shell")

local result = shell.run("unknown-command")
assert(result:exit_status() ~= 0)
assert(result:success() == false)
assert(result:failure() == true)
assert(result:stdout() == "")
assert(result:stderr() ~= "")
assert(result:combined_output() ~= "")
`
		err := L.DoString(code)
		assert.NoError(t, err)
	})
}
