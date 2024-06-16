local shell = require "xs.shell"
local debuglogger = require "xs.debuglogger"

local sshrc = function(override_config)
  local config = {
    -- default user home directory
    sshhome = os.getenv("HOME"),
  }
  if override_config then
    for k,v in pairs(override_config) do config[k] = v end
  end

  return function()
    -- check existing openssl command.
    if not shell.run("command -v openssl >/dev/null 2>&1") then
      error("sshrc requires openssl to be installed locally, but it's not. Aborting.")
    end

    debuglogger.printf("sshrc sshhome: %s", config.sshhome)

    local ret = shell.run([=[
function gen_sshrc_data() {
  local sshhome="]=] .. config.sshhome .. [=["
  if [ -f "$sshhome/.sshrc" ]; then
    local files=".sshrc"
    if [ -d $sshhome/.sshrc.d ]; then
      files="$files .sshrc.d"
    fi
    local total_file_size=$(tar cz -h -C $sshhome $files | wc -c)
    if [ $total_files_size -gt 65536 ]; then
      echo >&2 $'.sshrc.d and .sshrc files must be less than 64kb\ncurrent size: '$total_file_size' bytes'
      exit 1
    fi
    echo $(tar cz -h -C $sshhome $files | openssl enc -base64)
  else
    echo "No such file: $sshhome/.sshrc" >&2
    exit 1
  fi
}
gen_sshrc_data
    ]=])
    if not ret:success() then
      error(ret:combined_output())
    end

    sshrc_data = ret:stdout()

    if sshrc_data == nil or sshrc_data == "" then
      error("sshrc is empty")
    end

    -- return the script to be executed on the remote server
    return [=[
command -v openssl >/dev/null 2>&1 || { echo >&2 "sshrc requires openssl to be installed on the server, but it's not. Aborting."; exit 1; }
if [ -e /etc/motd ]; then cat /etc/motd; fi
if [ -e /etc/update-motd.d ]; then run-parts /etc/update-motd.d/ 2>/dev/null; fi
export SSHHOME=$(mktemp -d -t .$(whoami).sshrc.XXXX)
export SSHRCCLEANUP=$SSHHOME
trap "rm -rf $SSHRCCLEANUP; exit" 0

cat << 'EOF' > $SSHHOME/sshrc.bashrc
if [ -r /etc/profile ]; then source /etc/profile; fi
if [ -r ~/.bash_profile ]; then source ~/.bash_profile
elif [ -r ~/.bash_login ]; then source ~/.bash_login
elif [ -r ~/.profile ]; then source ~/.profile
fi
source $SSHHOME/.sshrc;
EOF

echo "]=] .. sshrc_data .. [=[" | tr -s ' ' $'\n' | openssl enc -base64 -d | tar mxz -C $SSHHOME

bash --rcfile $SSHHOME/sshrc.bashrc
exit $?
    ]=]
  end
end

return sshrc
