package cliclient

import (
	"errors"
	"fmt"

	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common/cliclient"
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Allows to manage personal data",
	Long:  `Allows to manage personal data`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loggedin, err := controller.ucauth.IsUserLoggedIn(cmd.Context())
		if err != nil {
			return err
		}

		if !loggedin {
			return fmt.Errorf("use command \"login\" before manage data")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.Println("Syncronizing...")
		err := controller.ucstore.Synchronize(cmd.Context())
		if err != nil {
			var serror *uccommon.ServerError
			if errors.As(err, &serror) {
				cmd.Printf("got server error: %v\n\n", serror.Err)
				return
			}

			cmd.Printf("Got error during syncronization: %s\n\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
