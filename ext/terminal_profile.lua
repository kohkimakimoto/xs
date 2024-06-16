--[[
terminal_profile is a hook that changes the macOS terminal profile.

Usage:
1. Put this source file into your ~/.xs directory.
2. Write the config like the following.
----
local terminal_profile = require "terminal_profile"

host "your-remote-server" {
  ssh_config= { ... }
  on_before_connect = { terminal_profile("Grass") }, -- change the terminal profile to "Grass"
  on_after_disconnect = { terminal_profile() }, -- change back to the original profile
}
----
--]]

local shell = require "xs.shell"

-- get the current terminal profile
function get_current_terminal_profile()
  local ret = shell.run([=[osascript -e 'tell application "Terminal" to return name of current settings of selected tab of front window']=])
  if not ret:success() then
    error(ret:combined_output())
  end
  return ret:stdout():gsub("\n$", "")
end

local terminal_profile = function(profile)
  if not profile then
    profile = get_current_terminal_profile()
  end
  return [=[osascript -e 'tell application "Terminal" to set current settings of first window to settings set "]=] .. profile .. [=["']=]
end

return terminal_profile
