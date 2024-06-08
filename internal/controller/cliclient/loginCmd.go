package cliclient

import (
	"github.com/spf13/cobra"
)

var (
	flagLoginLogin    string
	flagLoginPassword string
)

func init() {
	loginCmd.PersistentFlags().StringVar(&flagLoginLogin, "login", "", "set your login to authenticate")
	loginCmd.PersistentFlags().StringVar(&flagLoginPassword, "password", "", "set your password to authenticate")

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

		if flagLoginPassword == "" {
			wordPromptContent := promptContent{
				"Please provide a password.",
				"Enter a password:",
			}
			flagLoginPassword, err = promptPasswordInput(wordPromptContent)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		err = controller.ucauth.Login(ctx, flagLoginLogin, flagLoginPassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("login successfully")
	},
}
