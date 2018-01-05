# nova-cli

A convenient command line tool to pipe logs to [splunknova.com](https://www.splunknova.com) and search them.

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
This will install the `nova` binary to `$GOBIN`. If it isn't in your PATH, you can run `export PATH=$PATH:$GOBIN`

## Binaries

We may also have have binaries for download on the [latest release](https://github.com/splunknova/nova-cli/releases/latest).
Shout out on [slack](http://community.splunknova.com) if you need a particular binary!

# Usage

## Credentials

Get started by creating an account on [splunknova.com](https://www.splunknova.com/).

API Credentials can be conveniently saved in a `~/.nova` file by running

````
nova login
````

## Sending logs

You can pipe logs into nova by running 

````
echo "my first log" | nova

cat /var/log/system.log | nova

tail -f /var/log/system.log | nova
````

## Searching logs

Search for all lines containing the word "error"
````
nova search error
````

Only count the number of lines containing the word "error"
````
# shorthand
nova search error -c

# stats shorthand
nova search error -s count

# report command
nova search error -r "stats count"
````

Run stats aggregations and reporting on data
````
# SPL inspired syntax
nova search "my_key=" -r "stats count avg(my_key)"

# add transforms
nova search "bytes" -t "eval mb=gb*1024" -r "stats max(mb)"
````

## Sending Metrics

Create metric samples by running
````
# nova metric put <metric_name> <metric_value>
nova metric put cpu.usage 20

# tagging with dimensions
nova metric put cpu.usage 20 -d "region:us-east-1,role:webserver"
````

## Listing Metrics

````
nova metric ls
````

## Aggregating Metrics

````
# simple aggregations
nova metric get cpu.usage -a avg,max

# grouping by dimensions
nova metric get cpu.usage -a avg -g role
````


