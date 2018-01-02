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

If you haven't already, [sign up or log in][nova] to obtain your Splunk Nova API credentials to get started.

### Credentials

Save your Splunk Nova API client credentials in `~/.nova` file by running:

````
nova login
````
You will be prompted to enter your `Client ID` and `Client Secret`:
```
Please enter Client ID: <Your Client ID>
Please enter Client Secret: <Your Client Secret>
```
Once your credentials are entered, you should see:
```
INFO[0016] Login succeeded
```

## Send logs

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

You can also enter the `tail` command, followed by the file you’d like to view, which prints lines from the end of the file in reverse order:
```
tail -f /var/log/system.log | nova
```
Use the -f or --follow flag after tail, to show a real-time, streaming output of a changing file. It keeps watch and prints further data as it appears.

## Search logs

Search all logs containing the word "error"

```
nova search error
```
Returns:

```
count 0
```

Count only the number of lines containing the word "error"

```
nova search error -c
```
Returns
```
count 0
```

The `-s` or `stats` command calculates aggregate statistics over the
results set, such as average, count, and sum.

```
nova search error -s count
```

With the `stats` command you can specify a statistical function such as `count` to create a report of all errors. (How is this different than error -c?)

```
nova search error -r "stats count"
```

Run stats aggregations and reporting on data using Splunk Processing Language (SPL) inspired syntax. For example:

```
nova search "my_key=" -r "stats count avg(my_key)"
```
Returns a go routine that reports all usages of your Splunk Nova API credentials.

```
ERRO[0000] error communicating with splunknova. X-SPLUNK-REQ-ID=296890cdea349f739c8ffcd15828554d code:405, body:{"code":405,"message":"Error in 'search' command: Unable to parse the search: Comparator '=' is missing a term on the right hand side. in keywords"}
panic: runtime error: index out of range

goroutine 1 [running]:
github.com/splunknova/nova-cli/src.(*NovaSearch).Search(0xc420053c98, 0xc420014500, 0x19, 0x0, 0x0, 0x7fff5fbffb9f, 0x17, 0x0, 0x0, 0x0)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/src/search.go:77 +0xd09
github.com/splunknova/nova-cli/cmd.glob..func10(0x154cea0, 0xc42007d020, 0x1, 0x3)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/cmd/search.go:40 +0x3ac
github.com/spf13/cobra.(*Command).execute(0x154cea0, 0xc42007cf90, 0x3, 0x3, 0x154cea0, 0xc42007cf90)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:750 +0x2c1
github.com/spf13/cobra.(*Command).ExecuteC(0x154c120, 0xc420053f30, 0xc420053f38, 0x12c5a3e)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:831 +0x30e
github.com/spf13/cobra.(*Command).Execute(0x154c120, 0x0, 0x154bca0)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:784 +0x2b
github.com/splunknova/nova-cli/cmd.Execute()
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/cmd/root.go:52 +0x31
main.main()
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/nova.go:20 +0x20
```
Add transforming commands, a type of search command that orders the results into a data table

```
nova search "bytes" -t "eval mb=gb*1024" -r "stats max(mb)"
````
Returns

```
panic: runtime error: index out of range
goroutine 1 [running]:
github.com/splunknova/nova-cli/src.(*NovaSearch).Search(0xc420053c98, 0xc420014500, 0x17, 0x7fff5fbffb95, 0xf, 0x7fff5fbffba8, 0xd, 0x0, 0x0, 0x0)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/src/search.go:77 +0xd09
github.com/splunknova/nova-cli/cmd.glob..func10(0x154cea0, 0xc420086280, 0x1, 0x5)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/cmd/search.go:40 +0x3ac
github.com/spf13/cobra.(*Command).execute(0x154cea0, 0xc420086230, 0x5, 0x5, 0x154cea0, 0xc420086230)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:750 +0x2c1
github.com/spf13/cobra.(*Command).ExecuteC(0x154c120, 0xc420053f30, 0xc420053f38, 0x12c5a3e)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:831 +0x30e
github.com/spf13/cobra.(*Command).Execute(0x154c120, 0x0, 0x154bca0)
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/spf13/cobra/command.go:784 +0x2b
github.com/splunknova/nova-cli/cmd.Execute()
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/src/github.com/splunknova/nova-cli/cmd/root.go:52 +0x31
main.main()
	/private/tmp/nova-cli-20171218-11224-1uzn9g4/nova.go:20 +0x20
```

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
[novalogin]: https://www.splunknova.com/login

in Use:

You can post metrics to your datadog account by using:
