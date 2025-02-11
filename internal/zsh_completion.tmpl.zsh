# This is a zsh completion script for xs
# If you want to use this script, add the following line to your .zshrc
# ----------------------------------
# eval "$(xs zsh-completion)"
# ----------------------------------

_xs_hosts() {
  local -a __xs_hosts
  PRE_IFS=$IFS
  IFS=$'\n'
  __xs_hosts=($({{ .Executable }} zsh-completion --hosts | awk -F'\t' '{print $1":"$2}'))
  IFS=$PRE_IFS
  _describe -t host "host" __xs_hosts
}

_xs_builtin_commands() {
  local -a __xs_builtin_commands
  __xs_builtin_commands=(
    "list:List defined hosts"
    "ls:List defined hosts"
    "ssh-config:Output ssh_config to STDOUT"
    "zsh-completion:Output zsh completion script to STDOUT"
    "xscp-function:Output xscp function code to STDOUT"
  )
  _describe -t builtin_command "builtin command" __xs_builtin_commands
}

_xs() {
  local curcontext="$curcontext" state line expl ret=1
  local -a ssh_options
  typeset -A opt_args tsizes

  ssh_options=(
    '(-6)-4[force ssh to use IPv4 addresses only]'
    '(-4)-6[force ssh to use IPv6 addresses only]'
    '-A[enable forwarding of the authentication agent connection]'
    '-B+[bind to specified interface before attempting to connect]:interface:_net_interfaces'
    '-C[compress data]'
    '-D+[specify a dynamic port forwarding]:dynamic port forwarding:->dynforward'
    '-E+[append log output to file instead of stderr]:_files'
    '-F+[specify alternate config file]:config file:_files'
    '-G[output configuration and exit]'
    '-I+[specify smartcard device]:device:_files'
    '-J+[connect via a jump host]: :->userhost'
    '-K[enable GSSAPI-based authentication and forwarding]'
    '*-L+[specify local port forwarding]:local port forwarding:->forward'
    '-M[master mode for connection sharing]'
    '-m+[specify mac algorithms]: :->macs'
    "-N[don't execute a remote command]"
    '-O+[control an active connection multiplexing master process]:multiplex control command:((check\:"check master process is running" exit\:"request the master to exit" forward\:"request forward without command execution" stop\:"request the master to stop accepting further multiplexing requests" cancel\:"cancel existing forwardings with -L and/or -R" proxy))'
    '*-R+[specify remote port forwarding]:remote port forwarding:->forward'
    '-S+[specify location of control socket for connection sharing]:path to control socket:_files'
    '(-t)-T[disable pseudo-tty allocation]'
    '-V[show version number]'
    '-W+[forward standard input and output to host]:stdinout forward:->hostport'
    '(-x -Y)-X[enable (untrusted) X11 forwarding]'
    '(-x -X)-Y[enable trusted X11 forwarding]'
    '(-A)-a[disable forwarding of authentication agent connection]'
    '(-P)-b+[specify interface to transmit on]:bind address:_bind_addresses'
    '-c+[select encryption cipher]:encryption cipher:->ciphers'
    '-e+[set escape character]:escape character (or `none'\''):'
    '(-n)-f[go to background]'
    '-g[allow remote hosts to connect to local forwarded ports]'
    '*-i+[select identity file]:SSH identity file:_files -g "*(-.^AR)"'
    '-k[disable forwarding of GSSAPI credentials]'
    '-l+[specify login name]:login name:_ssh_users'
    '-m+[specify mac algorithms]: :->macs'
    '-n[redirect stdin from /dev/null]'
    '*-o+[specify extra options]:option string:->option'
    '-p+[specify port on remote host]:port number on remote host'
    '(-v)*-q[quiet operation]'
    '-s[invoke subsystem]'
    "(-T)*-t[force pseudo-tty allocation]"
    '(-q)*-v[verbose mode (multiple increase verbosity, up to 3)]'
    '-w+[request tunnel device forwarding]:local_tun[\:remote_tun] (integer or "any"):'
    '(-X -Y)-x[disable X11 forwarding]'
    '-y[send log info via syslog instead of stderr]'
  )

  _arguments -C -s \
    $ssh_options \
    '-h[show a help message and exit]' \
    ':builtin_command_or_destination:->builtin_command_or_destination' \
    '*::args:->command' \
    && ret=0

  case $state in
    builtin_command_or_destination)
      _xs_hosts
      _xs_builtin_commands
      return
      ;;

    # for ssh options
    option)
      if compset -P 1 '*='; then
        case "${IPREFIX#-o}" in
          (#i)(ciphers|macs|kexalgorithms|hostkeyalgorithms|pubkeyacceptedalgorithms)=)
          local sep
          zstyle -s ":completion:${curcontext}:" list-separator sep || sep=--
          if ! compset -P '[+-^]'; then
            _wanted prefix expl 'relative to default' compadd -S '' -d \
                "(
                   +\ $sep\ append\ to\ default\ list
                   -\ $sep\ remove\ from\ default\ list
                   ^\ $sep\ insert\ at\ head\ of\ default\ list
                )" - + - \^ && ret=0
          fi
          ;;
        esac
        case "${IPREFIX#-o}" in
        (#i)(batchmode|canonicalizefallbacklocal|checkhostip|clearallforwardings|compression|enableescapecommandline|enablesshkeysign|exitonforwardfailure|fallbacktorsh|forkafterauthentication|forward(agent|x11)|forwardx11trusted|gatewayports|gssapiauthentication|gssapidelegatecredentials|gssapikeyexchange|gssapirenewalforcesrekey|gssapitrustdns|hashknownhosts|hostbasedauthentication|identitiesonly|kbdinteractiveauthentication|tcpkeepalive|nohostauthenticationforlocalhost|passwordauthentication|permitlocalcommand|permitremoteopen|proxyusefdpass|stdinnull|streamlocalbindunlink|visualhostkey)=*)
          _wanted values expl 'truth value' compadd yes no && ret=0
          ;;
        (#i)addkeystoagent=*)
          _alternative \
            'timeouts: :_numbers -u seconds "time interval" :s:seconds m:minutes h:hours d:days w:weeks' \
            'values:value:(yes no ask confirm)' && ret=0
        ;;
        (#i)addressfamily=*)
          _wanted values expl 'address family' compadd any inet inet6 && ret=0
          ;;
        (#i)bindaddress=*)
          _wanted bind-addresses expl 'bind address' _bind_addresses && ret=0
          ;;
        (#i)bindinterface=*)
          _wanted bind-interfaces expl 'bind interface' _network_interfaces && ret=0
        ;;
        (#i)canonicaldomains=*)
          _message -e 'canonical domains (space separated)' && ret=0
          ;;
        (#i)canonicalizehostname=*)
          _wanted values expl 'truthish value' compadd yes no always && ret=0
          ;;
        (#i)canonicalizemaxdots=*)
          _message -e 'number of dots' && ret=0
          ;;
        (#i)canonicalizepermittedcnames=*)
          _message -e 'CNAME rule list (source_domain_list:target_domain_list, each pattern list comma separated)' && ret=0
          ;;
        (#i)ciphers=*)
          state=ciphers
          ;;
        (#i)certificatefile=*)
          _description files expl 'file'
          _files "$expl[@]" && ret=0
        ;;
        (#i)connectionattempts=*)
          _message -e 'connection attempts' && ret=0
          ;;
        (#i)connecttimeout=*)
          _numbers -u seconds timeout :s:seconds m:minutes h:hours d:days w:weeks && ret=0
          ;;
        (#i)controlmaster=*)
          _wanted values expl 'truthish value' compadd yes no auto ask autoask && ret=0
          ;;
        (#i)controlpath=*)
          _description files expl 'path to control socket'
          _files "$expl[@]" && ret=0
          ;;
        (#i)controlpersist=*)
          _alternative \
            'timeouts: :_numbers -u seconds timeout :s:seconds m:minutes h:hours d:days w:weeks' \
            'values:truth value:(yes no)' && ret=0
          ;;
        (#i)escapechar=*)
          _message -e 'escape character (or `none'\'')'
          ret=0
          ;;
        (#i)fingerprinthash=*)
          _values 'fingerprint hash algorithm' \
              md5 ripemd160 sha1 sha256 sha384 sha512 && ret=0
          ;;
        (#i)forwardx11timeout=*)
          _message -e 'timeout'
          ret=0
          ;;
        (#i)globalknownhostsfile=*)
          _description files expl 'global file with known hosts'
          _files "$expl[@]" && ret=0
          ;;
        (#i)hostname=*)
          _wanted hosts expl 'real host name to log into' _xs_hosts && ret=0
          ;;
        (#i)identityagent=*)
          _description files expl 'socket file'
          _files -g "*(-=)" "$expl[@]" && ret=0
        ;;
        (#i)identityfile=*)
          _description files expl 'SSH identity file'
          _files "$expl[@]" && ret=0
          ;;
        (#i)ignoreunknown=*)
          _message -e 'pattern list' && ret=0
          ;;
        (#i)ipqos=*)
          local descr
          if [[ $PREFIX = *\ *\ * ]]; then return 1; fi
          if compset -P '* '; then
            descr='QoS for non-interactive sessions'
          else
            descr='QoS [for interactive sessions if second value given, separated by white space]'
          fi
          _values $descr 'af11' 'af12' 'af13' 'af14' 'af22' \
              'af23' 'af31' 'af32' 'af33' 'af41' 'af42' 'af43' \
              'cs0' 'cs1' 'cs2' 'cs3' 'cs4' 'cs5' 'cs6' 'cs7' 'ef' \
              'lowdelay' 'throughput' 'reliability' && ret=0
          ;;
        (#i)(local|remote)forward=*)
          state=forward
          ;;
        (#i)dynamicforward=*)
          state=dynforward
          ;;
        (#i)kbdinteractivedevices=*)
          _values -s , 'keyboard-interactive authentication method' \
              'bsdauth' 'pam' 'skey' && ret=0
          ;;
        (#i)kexalgorithms=*)
          _wanted algorithms expl 'key exchange algorithm' _sequence compadd - \
              $(_call_program algorithms ssh -Q kex) && ret=0
          ;;
        (#i)gssapikexalgorithms=*)
          _wanted algorithms expl 'key exchange algorithm' _sequence compadd - \
              $(_call_program algorithms ssh -Q kex-gss) && ret=0
        ;;
        (#i)(local|knownhosts)command=*)
          _command_names -e && ret=0
          ;;
        (#i)loglevel=*)
          _values 'log level' QUIET FATAL ERROR INFO VERBOSE\
              DEBUG DEBUG1 DEBUG2 DEBUG3 && ret=0
          ;;
        (#i)macs=*)
          state=macs
          ;;
        (#i)numberofpasswordprompts=*)
          _message -e 'number of password prompts'
          ret=0
          ;;
        (#i)(pkcs11|securitykey)provider=*)
          _description files expl 'shared library'
          _files -g '*.(so|dylib)(|.<->)(-.)' "$expl[@]" && ret=0
          ;;
        (#i)port=*)
          _message -e 'port number on remote host'
          ret=0
          ;;
        (#i)preferredauthentications=*)
          _values -s , 'authentication method' gssapi-with-mic \
              hostbased publickey keyboard-interactive password && ret=0
          ;;
        (#i)proxyjump=*)
          compset -P "* "
          state=userhost
        ;;
        (#i)(hostkey|(hostbased|pubkey)accepted)algorithms=*)
	  _wanted key-types expl 'key type' _sequence compadd - \
              $(_call_program key-types ssh -Q key-sig) && ret=0
        ;;
        (#i)pubkeyauthentication=*)
          _wanted values expl 'enable' compadd yes no unbound host-bound && ret=0
        ;;
        (#i)protocol=*)
          _values -s , 'protocol version' \
              '1' \
              '2' && ret=0
          ;;
        (#i)(proxy|remote)command=*)
          _cmdstring && ret=0
          ;;
        (#i)rekeylimit=*)
          if compset -P "* "; then
            _numbers -u seconds "maximum time before renegotiating session key" \
                :s:seconds h:hours d:days w:weeks
          else
            _numbers -u bytes "maximum amount of data transmitted before renegotiating session key" \
                K:kilobytes M:megabytes G:gigabytes
          fi
          ret=0
          ;;
        (#i)requesttty=*)
          _values 'request a pseudo-tty' \
              'no[never request a TTY]' \
              'yes[always request a TTY when stdin is a TTY]' \
              'force[always request a TTY]' \
              'auto[request a TTY when opening a login session]' && ret=0
          ;;
        (#i)requiredrsasize=)
          _wanted sizes expl 'minimum size [1024]' compadd 1024 2048 4096 && ret=0
        ;;
        (#i)revokedhostkeys=*)
          _description files expl 'revoked host keys file'
          _files "$expl[@]" && ret=0
          ;;
        (#i)sendenv=*)
          _wanted envs expl 'environment variable' _parameters -g 'scalar*export*' && ret=0
          ;;
        (#i)serveralivecountmax=*)
          _message -e 'number of alive messages without replies before disconnecting'
          ret=0
          ;;
        (#i)serveraliveinterval=*)
          _message -e 'timeout in seconds since last data was received to send alive message'
          ret=0
          ;;
        (#i)streamlocalbindmask=*)
          _message -e 'octal mask' && ret=0
          ;;
        (#i)sessiontype=*)
          _wanted session-types expl "session type" compadd none subsystem default && ret=0
        ;;
        (#i)stricthostkeychecking=*)
          _wanted values expl 'value' compadd yes no ask accept-new off && ret=0
          ;;
        (#i)syslogfacility=*)
          _wanted facilities expl 'facility' compadd -M 'm:{a-z}={A-Z}' DAEMON USER AUTH LOCAL{0,1,2,3,4,5,6,7} && ret=0
          ;;
        (#i)(verifyhostkeydns|updatehostkeys)=*)
          _wanted values expl 'truthish value' compadd yes no ask && ret=0
          ;;
        (#i)transport=*)
          _values 'transport protocol' TCP SCTP && ret=0
          ;;
        (#i)tunnel=*)
          _values 'request device forwarding' \
              'yes' \
              'point-to-point' \
              'ethernet' \
              'no' && ret=0
          ;;
        (#i)tunneldevice=*)
          _message -e 'local_tun[:remote_tun] (integer or "any")'
          ret=0
          ;;
        (#i)userknownhostsfile=*)
          _description files expl 'user file with known hosts'
          _files "$expl[@]" && ret=0
          ;;
        (#i)user=*)
          _wanted users expl 'user to log in as' _ssh_users && ret=0
          ;;
        (#i)xauthlocation=*)
          _description files expl 'xauth program'
          _files "$expl[@]" -g '*(-*)' && ret=0
          ;;
        *) _message -e values value ;;
        esac
      else
        # Include, Host and Match not supported from the command-line
        # final GSSAPI options are not in upstream but are widely patched in
        _wanted values expl 'configure file option' \
            compadd -M 'm:{a-z}={A-Z} r:[^A-Z]||[A-Z]=* r:|=*' -q -S '=' - \
                AddKeysToAgent \
                AddressFamily \
                BatchMode \
                BindAddress \
                BindInterface \
                CanonicalDomains \
                CanonicalizeFallbackLocal \
                CanonicalizeHostname \
                CanonicalizeMaxDots \
                CanonicalizePermittedCNAMEs \
                CASignatureAlgorithms \
                CertificateFile \
                CheckHostIP \
                Ciphers \
                ClearAllForwardings \
                Compression \
                ConnectionAttempts \
                ConnectTimeout \
                ControlMaster \
                ControlPath \
                ControlPersist \
                DynamicForward \
                EnableEscapeCommandline \
                EnableSSHKeysign \
                EscapeChar \
                ExitOnForwardFailure \
                FingerprintHash \
                ForkAfterAuthentication \
                ForwardAgent \
                ForwardX11 \
                ForwardX11Timeout \
                ForwardX11Trusted \
                GatewayPorts \
                GlobalKnownHostsFile \
                GSSAPIAuthentication \
                GSSAPIDelegateCredentials \
                HashKnownHosts \
                HostbasedAcceptedAlgorithms \
                HostbasedAuthentication \
                HostKeyAlgorithms \
                HostKeyAlias \
                Hostname \
                IdentitiesOnly \
                IdentityAgent \
                IdentityFile \
                IgnoreUnknown \
                IPQoS \
                KbdInteractiveAuthentication \
                KbdInteractiveDevices \
                KexAlgorithms \
                KnownHostsCommand \
                LocalCommand \
                LocalForward \
                LogLevel \
                LogVerbose \
                MACs \
                NoHostAuthenticationForLocalhost \
                NumberOfPasswordPrompts \
                PasswordAuthentication \
                PermitLocalCommand \
                PermitRemoteOpen \
                PKCS11Provider \
                Port \
                PreferredAuthentications \
                ProxyCommand \
                ProxyJump \
                ProxyUseFdpass \
                PubkeyAuthentication \
                PubkeyAcceptedAlgorithms \
                RekeyLimit \
                RemoteCommand \
                RemoteForward \
                RequestTTY \
                RequiredRSASize \
                RevokedHostKeys \
                SecurityKeyProvider \
                SendEnv \
                ServerAliveCountMax \
                ServerAliveInterval \
                SetEnv \
                SessionType \
                StdinNull \
                StreamLocalBindMask \
                StreamLocalBindUnlink \
                StrictHostKeyChecking \
                SyslogFacility \
                TCPKeepAlive \
                Tunnel \
                TunnelDevice \
                UpdateHostKeys \
                User \
                UserKnownHostsFile \
                VerifyHostKeyDNS \
                VisualHostKey \
                XAuthLocation \
                GSSAPIClientIdentity \
                GSSAPIKeyExchange \
                GSSAPIRenewalForcesRekey \
                GSSAPIServerIdentity \
                GSSAPITrustDns \
                GSSAPIKexAlgorithms \
                && ret=0
      fi
      ;;
    forward)
      local port=false host=false listen=false bind=false
      if compset -P 1 '*:'; then
        if [[ $IPREFIX != (*=|)<-65535>: ]]; then
          if compset -P 1 '*:'; then
            if compset -P '*:'; then
              port=true
            else
              host=true
            fi
          else
            listen=true
            ret=0
          fi
        else
          if compset -P '*:'; then
            port=true
          else
            host=true
          fi
        fi
      else
        listen=true
        bind=true
      fi
      $port && { _message -e port-numbers 'port number'; ret=0 }
      $listen && { _message -e port-numbers 'listen-port number'; ret=0 }
      $host && { _wanted hosts expl host _xs_hosts && ret=0 }
      $bind && { _wanted bind-addresses expl bind-address _bind_addresses -S: && ret=0 }
      return ret
      ;;
    dynforward)
      _message -e port-numbers 'listen-port number'
      if ! compset -P '*:'; then
        _wanted bind-addresses expl bind-address _bind_addresses -qS:
      fi
      return 0
      ;;
    hostport)
      if compset -P '*:'; then
        _message -e port-numbers 'port number'
        ret=0
      else
        _wanted hosts expl host _xs_hosts && ret=0
      fi
      return ret
      ;;
    macs)
      _wanted macs expl 'MAC algorithm' _sequence compadd - $(_call_program macs ssh -Q mac)
      return
      ;;
    ciphers)
      _wanted ciphers expl 'encryption cipher' _sequence compadd - $(_call_program ciphers ssh -Q cipher)
      return
      ;;
    command)
      if (( $+opt_args[-s] )); then
        _wanted subsystems expl subsystem compadd sftp
        return
      fi
      local -a _comp_priv_prefix
      shift 1 words
      (( CURRENT-- ))
      _normal -p ssh
      return
      ;;
    userhost)
      _xs_hosts
      return 0
      ;;
  esac
}

compdef _xs xs
