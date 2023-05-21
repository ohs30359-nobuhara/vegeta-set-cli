package cmd

import (
	"github.com/spf13/cobra"
	"ohs30359/vegeta-cli/internal/scenario"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "setup",
	Short: "",
	Long:  ``,
	Run: func(c *cobra.Command, args []string) {
		path, e := c.PersistentFlags().GetString("scenario")
		if e != nil {
			panic(e.Error())
		}

		if e := scenario.Output(path); e != nil {
			panic(e.Error())
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("scenario", "s", "", "./xxx/xxx")
}
