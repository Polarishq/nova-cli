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

### Search for all lines containing the word "error"
````
nova search error
````

### Only count the number of lines containing the word "error"
````
# shorthand
nova search error -c

# stats shorthand
nova search error -s count

# report command
nova search error -r "stats count"
````

### Run stats aggregations and reporting on data
````
# SPL inspired syntax
nova search "my_key=" -r "stats count avg(my_key)"

# add transforms
nova search "bytes" -t "eval mb=gb*1024" -r "stats max(mb)"
````

## Sending Metrics

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

# grouping by dimensions (TODO, doesn't work yet)
nova metric get cpu.usage -a avg -g role
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

