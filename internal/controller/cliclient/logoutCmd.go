package cliclient

import "github.com/spf13/cobra"

func init() {
	root.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "The Logout command allows you to exit the system",
	Long:  `The Logout command allows you to exit the system. Only one client can use the system at a time. 
To use a different account, either set up another client or log out of the current one.`,
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
