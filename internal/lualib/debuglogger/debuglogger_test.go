package debuglogger

import (
	"bytes"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestPrintf(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	buf := new(bytes.Buffer)
	L.PreloadModule("debuglogger", Loader(debuglogger.New(buf, true, false)))

	code := `
local debuglogger = require("debuglogger")

debuglogger.printf("test %s", "message")
`
	err := L.DoString(code)
	assert.NoError(t, err)
	assert.Equal(t, "\x1b[2m[debug] test message\x1b[0m\n", buf.String())
}

func TestPrintfNoPrefix(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	buf := new(bytes.Buffer)
	L.PreloadModule("debuglogger", Loader(debuglogger.New(buf, true, false)))

	code := `
local debuglogger = require("debuglogger")

debuglogger.printf_no_prefix("test %s", "message")
`
	err := L.DoString(code)
	assert.NoError(t, err)
	assert.Equal(t, "\x1b[2mtest message\x1b[0m\n", buf.String())
}
