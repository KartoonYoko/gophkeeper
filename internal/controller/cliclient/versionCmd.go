package cliclient

import "github.com/spf13/cobra"

func init() {
	root.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Long:  `Print build version information`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln(controller.ucversion.Info())
	},
}
