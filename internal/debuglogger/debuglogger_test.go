package debuglogger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger_Printf(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, true, false)
		l.Printf("test %s", "message")
		assert.Equal(t, "\x1b[2m[debug] test message\x1b[0m\n", buf.String())
	})

	t.Run("no color", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, true, true)
		l.Printf("test %s", "message")
		assert.Equal(t, "[debug] test message\n", buf.String())
	})

	t.Run("no debug", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, false, false)
		l.Printf("test %s", "message")
		assert.Equal(t, "", buf.String())
	})
}

func TestLogger_PrintfNoPrefix(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, true, false)
		l.PrintfNoPrefix("test %s", "message")
		assert.Equal(t, "\x1b[2mtest message\x1b[0m\n", buf.String())
	})

	t.Run("no color", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, true, true)
		l.PrintfNoPrefix("test %s", "message")
		assert.Equal(t, "test message\n", buf.String())
	})

	t.Run("no debug", func(t *testing.T) {
		buf := new(bytes.Buffer)
		l := New(buf, false, false)
		l.PrintfNoPrefix("test %s", "message")
		assert.Equal(t, "", buf.String())
	})
}
