package cliclient

import "github.com/spf13/cobra"

func init() {
	root.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		err = controller.ucauth.Logout(ctx)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("login successfully")
	},
}
