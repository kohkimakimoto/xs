package internal

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"sort"
)

type Host struct {
	Name              string
	Description       string
	Hidden            bool
	SSHConfig         map[string]string
	OnBeforeConnect   []any
	OnAfterConnect    []any
	OnAfterDisconnect []any
}

func (h *Host) SortedSSHConfig() []map[string]string {
	values := make([]map[string]string, 0)
	names := make([]string, 0, len(h.SSHConfig))
	for name := range h.SSHConfig {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		v := h.SSHConfig[name]
		value := map[string]string{name: v}
		values = append(values, value)
	}

	return values
}

const luaHostTypeName = "Host*"

func registerLuaHostType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaHostTypeName)
	mt.RawSetString("__call", L.NewFunction(luaHostCall))
	mt.RawSetString("__index", L.NewFunction(luaHostIndex))
	mt.RawSetString("__newindex", L.NewFunction(luaHostNewIndex))
}

func xsHostFunc(L *lua.LState) int {
	l := L.GetTop()
	if l == 1 {
		// If it passes 1 argument, return host object for DSL style like `host "name" { description = "desc" })`
		name := L.CheckString(1)
		h, err := registerNewHost(L, name)
		if err != nil {
			L.RaiseError("failed to register new host: %v", err)
		}
		// return host object to Lua world
		L.Push(newLuaHost(L, h))
	} else if l == 2 {
		// If it passes 2 arguments, define a new host with function style like `host("name", { description = "desc" })`

		// host name and host config
		name := L.CheckString(1)
		tb := L.CheckTable(2)
		// register new host
		h, err := registerNewHost(L, name)
		if err != nil {
			L.RaiseError("failed to register new host: %v", err)
		}
		// apply host config
		tb.ForEach(func(k, v lua.LValue) {
			if key := lua.LVAsString(k); key != "" {
				if err := updateHost(h, key, v); err != nil {
					L.RaiseError("failed to parse host config: %v", err)
				}
			}
		})

		// return host object to Lua world
		L.Push(newLuaHost(L, h))
	} else {
		L.RaiseError("invalid number of arguments. want 1 or 2, got %d", l)
	}
	return 1
}

func registerNewHost(L *lua.LState, name string) (*Host, error) {
	// create new host object
	h := &Host{
		Name:        name,
		Description: "",
		SSHConfig:   map[string]string{},
	}

	// update config state
	cfg := getConfigFromLState(L)
	if err := cfg.AddHost(h); err != nil {
		return nil, err
	}
	return h, nil
}

func newLuaHost(L *lua.LState, host *Host) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = host
	L.SetMetatable(ud, L.GetTypeMetatable(luaHostTypeName))
	return ud
}

func updateHost(h *Host, key string, value lua.LValue) error {
	switch key {
	case "name":
		h.Name = lua.LVAsString(value)
	case "description":
		h.Description = lua.LVAsString(value)
	case "hidden":
		h.Hidden = lua.LVAsBool(value)
	case "ssh_config":
		if tb, ok := value.(*lua.LTable); ok {
			tb.ForEach(func(k, v lua.LValue) {
				if key := lua.LVAsString(k); key != "" {
					h.SSHConfig[key] = lua.LVAsString(v)
				}
			})
		} else {
			return fmt.Errorf("ssh_config must be a table but got %s", value.Type().String())
		}
	case "on_before_connect":
		// on_before_connect must be a table of string or function
		if tb, ok := value.(*lua.LTable); ok {
			hooks := make([]any, 0)
			tb.ForEach(func(_, v lua.LValue) {
				if vs, ok := v.(lua.LString); ok {
					// string value
					hooks = append(hooks, vs)
				} else if vfn, ok := v.(*lua.LFunction); ok {
					// function value
					hooks = append(hooks, vfn)
				}
			})
			h.OnBeforeConnect = hooks
		} else {
			return fmt.Errorf("on_before_connect must be a table but got %s", value.Type().String())
		}
	case "on_after_connect":
		// on_after_connect must be a table of string or function
		if tb, ok := value.(*lua.LTable); ok {
			hooks := make([]any, 0)
			tb.ForEach(func(_, v lua.LValue) {
				if vs, ok := v.(lua.LString); ok {
					// string value
					hooks = append(hooks, vs)
				} else if vfn, ok := v.(*lua.LFunction); ok {
					// function value
					hooks = append(hooks, vfn)
				}
			})
			h.OnAfterConnect = hooks
		} else {
			return fmt.Errorf("on_after_connect must be a table but got %s", value.Type().String())
		}
	case "on_after_disconnect":
		// on_after_disconnect must be a table of string or function
		if tb, ok := value.(*lua.LTable); ok {
			hooks := make([]any, 0)
			tb.ForEach(func(_, v lua.LValue) {
				if vs, ok := v.(lua.LString); ok {
					// string value
					hooks = append(hooks, vs)
				} else if vfn, ok := v.(*lua.LFunction); ok {
					// function value
					hooks = append(hooks, vfn)
				}
			})
			h.OnAfterDisconnect = hooks
		} else {
			return fmt.Errorf("on_after_disconnect must be a table but got %s", value.Type().String())
		}
	}
	return nil
}

func luaHostCall(L *lua.LState) int {
	h := checkHost(L)
	tb := L.CheckTable(2)
	// apply host config
	tb.ForEach(func(k, v lua.LValue) {
		if key := lua.LVAsString(k); key != "" {
			if err := updateHost(h, key, v); err != nil {
				L.RaiseError("failed to parse host config: %v", err)
			}
		}
	})
	L.Push(L.CheckUserData(1))
	return 1
}

func luaHostIndex(L *lua.LState) int {
	h := checkHost(L)
	key := L.CheckString(2)

	switch key {
	case "name":
		L.Push(lua.LString(h.Name))
		return 1
	case "description":
		L.Push(lua.LString(h.Description))
		return 1
	case "hidden":
		L.Push(lua.LBool(h.Hidden))
		return 1
	case "ssh_config":
		tb := L.NewTable()
		for k, v := range h.SSHConfig {
			tb.RawSetString(k, lua.LString(v))
		}
		L.Push(tb)
		return 1
	case "on_before_connect":
		tb := L.NewTable()
		for i, v := range h.OnBeforeConnect {
			switch vv := v.(type) {
			case lua.LString:
				tb.RawSetInt(i+1, vv)
			case *lua.LFunction:
				tb.RawSetInt(i+1, vv)
			}
		}
		L.Push(tb)
		return 1
	case "on_after_connect":
		tb := L.NewTable()
		for i, v := range h.OnAfterConnect {
			switch vv := v.(type) {
			case lua.LString:
				tb.RawSetInt(i+1, vv)
			case *lua.LFunction:
				tb.RawSetInt(i+1, vv)
			}
		}
		L.Push(tb)
		return 1
	case "on_after_disconnect":
		tb := L.NewTable()
		for i, v := range h.OnAfterDisconnect {
			switch vv := v.(type) {
			case lua.LString:
				tb.RawSetInt(i+1, vv)
			case *lua.LFunction:
				tb.RawSetInt(i+1, vv)
			}
		}
		L.Push(tb)
		return 1
	default:
		L.Push(lua.LNil)
		return 1
	}
}

func luaHostNewIndex(L *lua.LState) int {
	h := checkHost(L)
	key := L.CheckString(2)
	value := L.CheckAny(3)

	if err := updateHost(h, key, value); err != nil {
		L.RaiseError("failed to parse host config: %v", err)
	}
	return 0
}

func checkHost(L *lua.LState) *Host {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Host); ok {
		return v
	}
	L.ArgError(1, "Host object expected")
	return nil
}
