--[[
dhcpd_leases module parses the /var/db/dhcpd_leases file and provides functions to search for specific entries.

Usage:
1. Put this source file into your ~/.xs directory.
2. Write the config like the following.
----
local dhcpd_leases = require "dhcpd_leases"
local leases = dhcpd_leases.parse()

-- Find all entries with a specific key-value pair
local results = leases:find("name", "utm-ubuntu")

-- Find the first entry with a specific key-value pair
local result = leases:findOne("name", "utm-ubuntu")
----
--]]

local dhcpd_leases = {}

local function parse_entry(block)
  local entry = {}
  for line in block:gmatch("[^\n]+") do
    local key, value = line:match("([%w_]+)%s*=%s*(.+)")
    if key and value then
      entry[key] = value
    end
  end
  return entry
end

dhcpd_leases.parse = function(file_path)
  file_path = file_path or "/var/db/dhcpd_leases"
  local leases = {}

  -- open file
  local file = io.open(file_path, "r")
  if not file then
    return leases
  end

  local file_content = file:read("*all")
  file:close()

  -- parse entries
  for block in file_content:gmatch("{.-}") do
    table.insert(leases, parse_entry(block))
  end

  -- add find method
  leases.find = function(self, key, value)
    local results = {}
    for _, lease in ipairs(self) do
      if lease[key] == value then
        table.insert(results, lease)
      end
    end
    return results
  end

  -- add findOne method
  leases.findOne = function(self, key, value)
    for _, lease in ipairs(self) do
      if lease[key] == value then
        return lease
      end
    end
    return nil
  end

  return leases
end

return dhcpd_leases
