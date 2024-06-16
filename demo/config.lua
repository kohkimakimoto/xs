-- add "ext" diretory to package.path
package.path = package.path .. ";" .. xs.config_dir .. "/../ext/?.lua"

local terminal_profile = require "terminal_profile"
local sshrc = require "sshrc"

host "demo-dev-server1" {
  description = "This host for changing color demo",
  ssh_config = {
    HostName = "127.0.0.1",
    Port = "2222",
    User = "xs-test-user",
    StrictHostKeyChecking = "no",
    IdentityFile = xs.config_dir .. "/key",
  },
  on_before_connect = { terminal_profile("Pro Red") }, -- You need to create a "Pro Red" profile in the terminal preferences.
  on_after_disconnect = { terminal_profile() },
}

host "demo-dev-server2" {
  description = "This host for sshrc hook demo",
  ssh_config = {
    HostName = "127.0.0.1",
    Port = "2222",
    User = "xs-test-user",
    StrictHostKeyChecking = "no",
    IdentityFile = xs.config_dir .. "/key",
  },
  on_after_connect = { sshrc({ sshhome = xs.config_dir}) }
}

host "demo-production-server1" {
  description = "demo production server1",
  ssh_config = {
    HostName = "127.0.0.1",
    Port = "2222",
    User = "xs-test-user",
    StrictHostKeyChecking = "no",
    IdentityFile = xs.config_dir .. "/key",
  }
}
