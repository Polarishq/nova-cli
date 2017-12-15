# nova-cli

A convenient command line tool to pipe logs to [splunknova.com](https://www.splunknova.com) and search them.

# Usage

## Credentials
Credentials can be conveniently saved in a `~/.nova` file by running

````
nova login
````

## Sending logs
````
echo "my first log" | nova

cat /var/log/system.log | nova

tail -f /var/log/system.log | nova
````

## Searching logs
````
nova search ERROR

nova search ERROR -c

nova search ERROR -s count

nova search "my_key=" -r "stats count avg(my_key)"
````

# Installation

## macOS

````
brew tap splunknova/nova-cli
brew install nova-cli
````

## Linux & Windows

For now, you'll need to ensure `go` is installed and `GOROOT`, `GOPATH`, and `GOBIN` are set.

_We need help with making Linux and Windows installations better, please send a PR!_

````
go get github.com/splunknova/nova-cli
cd $GOPATH/src/github.com/splunknova/nova-cli
go install nova.go
````
This will install the `nova` binary to `$GOBIN`

