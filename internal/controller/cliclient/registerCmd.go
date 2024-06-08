package cliclient

import (
	"github.com/spf13/cobra"
)

var (
	flagRegisterLogin    string
	flagRegisterPassword string
)

func init() {
	registerCmd.Flags().StringVar(&flagRegisterLogin, "login", "", "set your login to authenticate")
	registerCmd.Flags().StringVar(&flagRegisterPassword, "password", "", "set your password to authenticate")

	root.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
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

		err = cmd.ValidateRequiredFlags()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if flagRegisterPassword == "" {
			wordPromptContent := promptContent{
				"Please provide a password.",
				"Enter a password:",
			}
			flagRegisterPassword, err = promptPasswordInput(wordPromptContent)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		err = controller.ucauth.Register(ctx, flagRegisterLogin, flagRegisterPassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("register successfully")
	},
}
