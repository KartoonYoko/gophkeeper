package cliclient

import (
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// DO stuff here
	},
}
