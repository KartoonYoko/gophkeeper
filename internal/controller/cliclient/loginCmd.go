package cliclient

import (
	"github.com/spf13/cobra"
)

func init() {
	loginCmd.Flags().String("login", "", "set your login to login")
	loginCmd.Flags().String("password", "", "set your password to login")

	root.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		err = cmd.MarkFlagRequired("login")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		flogin, err := validateRequiredStringFlag(cmd, "login")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		fpassword, err := validateRequiredStringFlag(cmd, "login")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		if fpassword == "" {
			wordPromptContent := promptContent{
				"Please provide a password.",
				"Enter a password:",
			}
			fpassword, err = promptPasswordInput(wordPromptContent)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		err = controller.ucauth.Login(ctx, flogin, fpassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("login successfully")
	},
}
