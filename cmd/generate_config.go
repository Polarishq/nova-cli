package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "generate-config",
	Short: "A utility for generating config files with your credentials",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

type configTemplateData struct {
	AuthHeader string
	NovaURL    string
}

func init() {
	rootCmd.AddCommand(generateConfigCmd)
}
