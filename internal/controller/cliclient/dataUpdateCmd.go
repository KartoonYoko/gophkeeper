package cliclient

import (
	"errors"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	"github.com/spf13/cobra"
)

func init() {
	dataUpdateCmd.Flags().String("dataid", "", "ID of data to update")

	dataCmd.AddCommand(dataUpdateCmd)
}

var dataUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update private data by id",
	Long:  `It allows you to update data`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		dataid, err := cmd.Flags().GetString("dataid")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if dataid == "" {
			pc := promptContent{
				errorMsg: "Please, enter ID of data",
				label:    "Enter data ID",
			}
			dataid, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		pc := promptContent{
			errorMsg: "Please, enter text",
			label:    "Enter new text",
		}
		text, err := promptTextInput(pc)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// TODO get data type and by type decide what method to use
		err = controller.ucstore.UpdateTextData(ctx, dataid, text)
		if err != nil {
			var serror *clientstore.ServerError
			if errors.As(err, &serror) {
				cmd.Printf("Successfull updated locally, but got server error: %s. ", serror.Err)
				return
			}
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Data successfully updated")
	},
}
