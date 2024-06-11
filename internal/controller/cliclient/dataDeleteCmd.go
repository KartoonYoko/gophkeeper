package cliclient

import (
	"errors"

	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common/cliclient"
	"github.com/spf13/cobra"
)

func init() {
	dataDeleteCmd.Flags().String("dataid", "", "ID of data to delete")

	dataCmd.AddCommand(dataDeleteCmd)
}

var dataDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete private data by id",
	Long:  `It allows you to delete data`,
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
				errorMsg: "Please, enter ID of data to delete",
				label:    "Enter data ID",
			}
			dataid, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		err = controller.ucstore.RemoveDataByID(ctx, dataid)
		if err != nil {
			var serror *uccommon.ServerError
			if errors.As(err, &serror) {
				cmd.Printf("Successfull deleted locally, but got server error: %s. ", serror.Err)
				return
			}
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Data successfully deleted")
	},
}
