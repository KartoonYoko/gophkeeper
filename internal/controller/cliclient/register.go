package cliclient

import (
	"github.com/spf13/cobra"
)

var (
	// флаги
	flagRegisterLogin    string
	flagRegisterPassword string
)

func init() {
	registerCmd.PersistentFlags().StringVar(&flagRegisterLogin, "login", "", "set your login to authenticate")
	registerCmd.PersistentFlags().StringVar(&flagRegisterPassword, "password", "", "set your password to authenticate")

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

		err = controller.ucauth.Register(ctx, flagRegisterLogin, flagRegisterPassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("register successfully")
	},
}
