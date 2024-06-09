package cliclient

import "github.com/spf13/cobra"

func init() {
	dataGetCmd.AddCommand(dataGetListCmd)
}

var dataGetListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// var err error
		// ctx := cmd.Context()

		// todo usecase get list
	},
}
