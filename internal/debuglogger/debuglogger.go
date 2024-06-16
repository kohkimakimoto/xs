package debuglogger

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
	"io"
)

type Logger struct {
	Writer  io.Writer
	IsDebug bool
	NoColor bool
}

func New(w io.Writer, isDebug bool, noColor bool) *Logger {
	return &Logger{
		Writer:  w,
		IsDebug: isDebug,
		NoColor: noColor,
	}
}

func (l *Logger) Printf(format string, a ...any) {
	if l.IsDebug {
		var txt string
		if l.NoColor {
			txt = fmt.Sprintf("[debug] "+format, a...)
		} else {
			txt = fmt.Sprintf(text.Faint.Sprint("[debug] "+format), a...)
		}
		if len(txt) == 0 || txt[len(txt)-1] != '\n' {
			txt += "\n"
		}
		_, _ = fmt.Fprint(l.Writer, txt)
	}
}

func (l *Logger) PrintfNoPrefix(format string, a ...any) {
	if l.IsDebug {
		var txt string
		if l.NoColor {
			txt = fmt.Sprintf(format, a...)
		} else {
			txt = fmt.Sprintf(text.Faint.Sprint(format), a...)
		}
		if len(txt) == 0 || txt[len(txt)-1] != '\n' {
			txt += "\n"
		}
		_, _ = fmt.Fprint(l.Writer, txt)
	}
}

func Bind(app *cli.App, l *Logger) {
	app.Metadata["logger"] = l
}

func Get(cCtx *cli.Context) *Logger {
	l, ok := cCtx.App.Metadata["logger"].(*Logger)
	if !ok {
		return nil
	}
	return l
}
