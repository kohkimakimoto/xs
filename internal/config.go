package internal

import (
	"fmt"
	"github.com/kohkimakimoto/xs/internal/debuglogger"
	luadebuglogger "github.com/kohkimakimoto/xs/internal/lualib/debuglogger"
	"github.com/kohkimakimoto/xs/internal/lualib/shell"
	"github.com/kohkimakimoto/xs/internal/lualib/template"
	"github.com/urfave/cli/v2"
	"github.com/yuin/gopher-lua"
	"os"
	"path/filepath"
)

type Config struct {
	Filepath    string
	Hosts       []*Host
	DebugLogger *debuglogger.Logger
}

func (cfg *Config) NewHostFilter() *HostFilter {
	return &HostFilter{
		hosts: cfg.Hosts,
	}
}

func (cfg *Config) AddHost(h *Host) error {
	for _, host := range cfg.Hosts {
		if host.Name == h.Name {
			return fmt.Errorf("host %s already registered", h.Name)
		}
	}
	cfg.Hosts = append(cfg.Hosts, h)
	return nil
}

const LuaConfigKey = "*__xs_config"

// registerConfig registers the config object in the Lua state.
func registerConfig(L *lua.LState) {
	config := &Config{
		Hosts: []*Host{},
	}
	ud := L.NewUserData()
	ud.Value = config
	L.SetGlobal(LuaConfigKey, ud)
}

func getConfigFromLState(L *lua.LState) *Config {
	v := L.GetGlobal(LuaConfigKey)
	if v == nil {
		panic("does not register config object in the Lua state")
	}
	ud, ok := v.(*lua.LUserData)
	if !ok {
		panic("detects invalid config object in the Lua state")
	}
	cfg, ok := ud.Value.(*Config)
	if !ok {
		panic("detects invalid config object in the Lua state")
	}
	return cfg
}

func newLState() *lua.LState {
	L := lua.NewState()

	// define built-in functions
	L.SetGlobal("host", L.NewFunction(xsHostFunc))

	// register config object
	registerConfig(L)

	// register types
	registerLuaHostType(L)

	return L
}

type ConfigLoadError struct {
	Err  error
	Path string
}

func (e *ConfigLoadError) Error() string {
	return "failed to load config (" + e.Err.Error() + ")"
}

func newConfig(cCtx *cli.Context) (*Config, *lua.LState, error) {
	configFilePath := getConfigFilePath()
	if absPath, err := filepath.Abs(configFilePath); err == nil {
		configFilePath = absPath
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, &ConfigLoadError{Err: fmt.Errorf("no config file: %s", configFilePath), Path: configFilePath}
	}

	// Create Lua state that corresponds to the config.
	L := newLState()

	// Get config object from Lua state
	cfg := getConfigFromLState(L)
	// Setup config object
	cfg.Filepath = configFilePath

	// Register "xs" predefined global variable
	xsObject := L.NewTable()
	L.SetGlobal("xs", xsObject)
	xsObject.RawSetString("config_file", lua.LString(configFilePath))
	xsObject.RawSetString("config_dir", lua.LString(filepath.Dir(configFilePath)))

	// Load built-in modules
	L.PreloadModule("xs.debuglogger", luadebuglogger.Loader(debuglogger.Get(cCtx)))
	L.PreloadModule("xs.shell", shell.Loader)
	L.PreloadModule("xs.template", template.Loader)

	// Extend package.path
	// The directory of the config file is added to the package.path.
	dir := filepath.Dir(configFilePath)
	var additionalPath string
	if os.PathSeparator == '/' { // unix-like
		additionalPath = dir + "/?.lua;"
	} else {
		additionalPath = dir + "\\?.lua;"
	}

	if err := L.DoString(`package.path = "` + additionalPath + `" .. package.path`); err != nil {
		return nil, nil, err
	}

	// Load config file
	if err := L.DoFile(configFilePath); err != nil {
		return nil, nil, &ConfigLoadError{Err: err, Path: configFilePath}
	}

	return cfg, L, nil
}
