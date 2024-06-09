package cliclient

import "github.com/spf13/cobra"

func init() {
	dataCmd.AddCommand(dataGetCmd)
}

var dataGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
