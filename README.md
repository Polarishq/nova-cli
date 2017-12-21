# Nova CLI

A convenient command-line tool for sending and searching logs using [splunknova.com](https://www.splunknova.com).

## Get Started
Get started by creating an account on https://www.splunknova.com.

## Install

Set up your [Go] environment.  Refer to the official Go documentation for more details: https://golang.org/doc/code.html. Once Go is downloaded and installed, you'll need to set your `GOROOT`, `GOPATH`, and `GOBIN`.

### macOS

[Homebrew] is a package manager for macOS that makes it easy to install Nova CLI. In a terminal window, run:

```
brew tap splunknova/nova-cli
```

```
brew install nova-cli
```

### Linux & Windows

**Linux**: By default Go is installed to directory `/usr/local/go/`, and the `GOROOT` environment variable is set to `/usr/local/go/bin`.

**Windows**: By default Go is installed in the directory `C:\Go`, the `GOROOT` environment variable is set to `C:\Go\`, and the bin directory is added to your Path (`C:\Go\bin`).

To install Nova CLI using Go, in the command-line run:

```
go get github.com/splunknova/nova-cli
```
Change directories into the `nova-cli` repository.
```
cd $GOPATH/src/github.com/splunknova/nova-cli
```
Install the `nova` binary to `$GOBIN`.
```
go install nova.go
```
 If it isn't in your PATH, you can run `export PATH=$PATH:$GOBIN`


## Usage

If you haven't already, [sign up or log in][nova]  to obtain client credentials to get started.

### Credentials

Save your Splunk Nova client credentials in `~/.nova` file by running:

````
nova login
````
You will be prompted to enter your `Client ID` and `Client Secret`:
```
Please enter Client ID: <Your Client ID>
Please enter Client Secret: <Your Client Secret>
```
Once your credentials are entered, you
```
INFO[0016] Login succeeded
```

### Send logs

You can pipe logs into Splunk Nova by running:

```
echo "my first log" | nova
```
This sends a log string: `"my first log"` to nova. You can then search your logs from the CLI using `nova search`. For example:

```
nova search "my first log"
```
returns a list of `my first logs` sent to nova:

```
2018-1-19T19:35:01.000+00:00 my first log
2018-1-18T23:53:08.000+00:00 my first log
2018-1-18T23:52:38.000+00:00 my first cli log
```
One example of a `cat ` command for system log files would be to pipe `system.log` to Splunk Nova:
```
cat /var/log/system.log | nova
```
and then search:

```
nova search system.log
```

Returns:
```
2017-12-21T00:02:07.000+00:00 	ASL Module "com.apple.authd" sharing output destination "/var/log/system.log" with ASL Module "com.apple.asl".
2017-12-21T00:02:07.000+00:00 	ASL Module "com.apple.authd" sharing output destination "/var/log/system.log" with ASL Module "com.apple.asl".
2017-12-18T23:54:15.000+00:00 	ASL Module "com.apple.authd" sharing output destination "/var/log/system.log" with ASL Module "com.apple.asl".
```

When you need to look for specific events from log files, your may want to use 'tail'
```
tail -f /var/log/system.log | nova
````
The `-f` follow flag causes  tail, after printing lines from the end of the file in reverse  order, to keep watch and print further data as it appears. Super helpful for ...? (One last tip about the tail command. Well, when you review the help page for tail, read the -s, or sleep, option as it automates the review of log files, which is a good security precaution.) transform this method to "tail" multiple files since "tail" is able to watch more than one file...Since your node might stay alive for many days, you may need some logrotation tool to replace the file you're lurking?

### Search logs

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


[Go]: https://golang.org/dl/
[homebrew]: https://brew.sh/
[nova]: https://www.splunknova.com/


in Use:

You can post metrics to your datadog account by using:
