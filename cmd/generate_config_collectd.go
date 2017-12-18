package cmd

import (
	"github.com/spf13/cobra"
	"text/template"
	"os"
)


// collectdCmd represents the collectd command
var collectdCmd = &cobra.Command{
	Use:   "collectd",
	Short: "Show a sample collectd configuration",
	PreRun: Authorize,
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, _ := template.New("Foo").Parse(collectdTemplateStr)
		f := configTemplateData{AuthHeader: AuthHeader,NovaURL: NovaURL}
		tmpl.Execute(os.Stdout, f)
	},
}

func init() {
	generateConfigCmd.AddCommand(collectdCmd)
}

var collectdTemplateStr = `
# Turn on the write_http plugin in collectd.conf
# Typically found at /usr/local/etc/collectd.conf

LoadPlugin write_http
<Plugin write_http>
        <Node "splunknova">
                URL "{{.NovaURL}}/v1/metrics?type=collectd"
                Header "Authorization: {{.AuthHeader}}"
                Format "JSON"
                Metrics true
                Notifications false
                StoreRates false
                BufferSize 4096
                LowSpeedLimit 0
                Timeout 0
                LogHttpError true
        </Node>
</Plugin>

`