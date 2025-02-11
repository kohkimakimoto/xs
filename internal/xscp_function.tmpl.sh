# This is a shell function definition for xs
# If you want to use this script, add the following line to your .zshrc
# ----------------------------------
# eval "$(xs xscp-function)"
# ----------------------------------

function {{ .Name }}() {
  local tmp_ssh_config=$(mktemp)
  cleanup() {
    rm -f "$1"
  }
  trap "cleanup \"$tmp_ssh_config\"" EXIT
  {{ .Executable }} ssh-config > "$tmp_ssh_config"
  scp -F "$tmp_ssh_config" "$@"
}

