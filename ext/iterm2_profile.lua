--[[
iterm2_profile is a hook that changes the iTerm2 profile.

Usage:
1. Put this source file into your ~/.xs directory.
2. Write the config like the following.
----
local iterm2_profile = require "iterm2_profile"

host "your-remote-server" {
  ssh_config= { ... }
  on_before_connect = { iterm2_profile("Remote") },
  on_after_disconnect = { iterm2_profile("Default") },
}
----
--]]
local iterm2_profile = function(profile)
  return [=[printf "\033]1337;SetProfile=]=] .. profile .. [=[\a"]=]
end

return iterm2_profile
