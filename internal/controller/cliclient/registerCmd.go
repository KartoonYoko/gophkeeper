package cliclient

import (
	"github.com/spf13/cobra"
)

func init() {
	registerCmd.Flags().String("login", "", "set your login to authenticate")
	registerCmd.Flags().String("password", "", "set your password to authenticate")

	root.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates a new user account",
	Long:  `The registration command creates a new user account`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		flogin, err := cmd.Flags().GetString("login")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		if flogin == "" {
			pc := promptContent{
				errorMsg: "Please, enter your login",
				label:    "Enter your login",
			}
			flogin, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		fpassword, err := cmd.Flags().GetString("password")
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

		err = controller.ucauth.Register(ctx, flogin, fpassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("register successfully")
	},
}
