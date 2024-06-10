package cliclient

import "github.com/spf13/cobra"

func init() {
	root.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Allows to manage personal data",
	Long:  `Allows to manage personal data`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
