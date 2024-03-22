# Overview

Tincl (Telnet INteractive CLient) is an interactive telnet client with history and automation via scripting. Tincl supports only text based telnet sessions and can't process binary data. All incoming and outgoing data is converting to strings.

# License
Tincl is released under the GPL 3.0 license. See [LICENSE.txt](LICENSE.txt)

# Current state

At the moment tincl has alfa quality code. There is no timeouts, test cases and fuzzing. You should not use it with untrusted telnet servers.

# Usage

There is two modes, interactive and non interactive. By default is using non interactive mode. This mean that tincl is reading stdin without readline. In interactive mode (-i option), will be enable ishell instance with readline support. Currently in interactive mode you should use 'crlf' command for sending pure '\r\n' or '\n' to telnet server, because ishell filters '\n'. Tincl supports lua scripts via -s option in CLI and via 'script' command inside interactive mode. But only -s currently supports telnet session processing.

Example:

```
>tincl -H smtp.google.com -P 25 -n -i -s test.lua
* Connecting to smtp.google.com:25
* Successful connected to smtp.google.com:25
* SMTP Greeting: 	220 mx.google.com ESMTP b9-20020ac24109000000b004fe2ad08fc3si1139642lfi.217 - gsmtp
* SMTP Answer after helo: 	250 mx.google.com at your service
>>> mail from: zzz
555 5.5.2 Syntax error. b9-20020ac24109000000b004fe2ad08fc3si1139642lfi.217 - gsmtp
>>>
```

CLI Options:
```
Usage of tincl:

  -n, --disable-greeting   Do not read greeting from telnet server.
  -H, --host string        Connection host
  -i, --interactive        Enable interactive mode
  -P, --port int           Connection port (default 23)
  -s, --script string      Script for execution
  -t, --tls                Use TLS mode
```

Interactive commands:
```
>>> help

Commands:
  clear       clear the screen
  crlf        Send crlf
  exit        exit the program
  help        display help
  script      Run lua script
```

# TODO

* Disable main read loop after loading lua script inside interactive mode.
* Allow to use '\n' only in interactive mode (require PR to ishell).

# Thanks

[viper](github.com/spf13/viper/) - is a complete configuration solution for Go applications including 12-Factor apps from spf13 team.
[gopher-lua](https://github.com/yuin/gopher-lua) - is a Lua5.1(+ goto statement in Lua5.2) VM and compiler written in Go.
[go-telnet](https://github.com/reiver/go-telnet) - Thanks for the golang telnet library.
[ishell](https://github.com/abiosoft/ishell/) - is an interactive shell library for creating interactive cli applications.
