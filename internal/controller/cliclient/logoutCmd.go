package cliclient

import "github.com/spf13/cobra"

func init() {
	logoutCmd.Flags().Bool("force", false, "force logout even if server not responding")

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

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if force {
			err = controller.ucauth.LogoutForce(ctx)
		} else {
			err = controller.ucauth.Logout(ctx)
		}
		
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("login successfully")
	},
}
