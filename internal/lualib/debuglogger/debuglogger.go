package debuglogger

import (
	"fmt"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	"github.com/yuin/gopher-lua"
)

func Loader(l *debuglogger.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		tb := L.NewTable()
		L.SetFuncs(tb, map[string]lua.LGFunction{
			"printf":           printf(l),
			"printf_no_prefix": printfNoPrefix(l),
		})
		L.Push(tb)
		return 1
	}
}

func printf(l *debuglogger.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		top := L.GetTop()
		format := ""
		values := make([]any, 0, top-1)
		for i := 1; i <= top; i++ {
			if i == 1 {
				format = L.CheckString(i)
				continue
			}
			values = append(values, toGoValue(L.Get(i)))
		}
		l.Printf(format, values...)
		return 0
	}
}

func printfNoPrefix(l *debuglogger.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		top := L.GetTop()
		format := ""
		values := make([]any, 0, top-1)
		for i := 1; i <= top; i++ {
			if i == 1 {
				format = L.CheckString(i)
				continue
			}
			values = append(values, toGoValue(L.Get(i)))
		}
		l.PrintfNoPrefix(format, values...)
		return 0
	}
}

func toGoValue(lv lua.LValue) any {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]any)
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(toGoValue(key))
				ret[keystr] = toGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]any, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, toGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}
