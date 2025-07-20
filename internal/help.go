package internal

const rootHelpTemplate = `Usage: xs [options] builtin_command|destination [command [args ...]]

XS is a SSH command wrapper that enhances your SSH operations.

Options:
   -h, --help     Show this help message and exit
   You can also use ssh command options. Check 'man ssh' for more information.

Builtin commands:{{template "visibleCommandCategoryTemplate" .}}

Destination Hosts:
   You can define destination hosts in the configuration file.

Environment variables:
   XS_CONFIG_FILE  Path to the configuration file. Default is ~/.xs/config.lua
   XS_DEBUG        If set to "true", XS will output debug information.
   XS_NO_COLOR     If set to "true", XS will not output color codes in debug information.

Version: {{ .Version }}
Commit: {{ .Metadata.CommitHash }}
{{template "copyrightTemplate" .}}
`

const helpTemplate = `Usage: {{template "usageTemplate" .}}

{{template "helpNameTemplate" .}}

Options:{{template "visibleFlagTemplate" .}}
`
