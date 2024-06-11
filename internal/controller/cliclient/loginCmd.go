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
	Short: "Login to gophkeeper",
	Long: `The Login command allows you to authenticate. 
	Only authenticated users can store personal information.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		// ввод логина
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
		
		// ввод пароля
		fpassword, err := cmd.Flags().GetString("login")
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
