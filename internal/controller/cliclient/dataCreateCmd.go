package cliclient

import (
	"errors"
	"strings"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	"github.com/spf13/cobra"
)

func init() {
	dataCreateCmd.Flags().String("datatype", "TEXT", "Set data type that you want create. Possible values: TEXT, CREDENTIALS, BANK_CARD, BINARY.")
	dataCmd.AddCommand(dataCreateCmd)
}

type dataType string

var dataCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create private data",
	Long:  `It allows you to create personal data. By default, it creates text data.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		datatype, err := cmd.Flags().GetString("datatype")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		dt := dataType(datatype)
		if !dt.isValid() {
			cmd.PrintErrln("Invalid data type. Please, use TEXT, CREDENTIALS, BANK_CARD, BINARY.")
			return
		}

		if dt == "TEXT" {
			pc := promptContent{
				errorMsg: "Please, enter text",
				label:    "Enter text",
			}
			text, err := promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = controller.ucstore.CreateTextData(ctx, text)
			if err != nil {
				var serror *clientstore.ServerError
				if errors.As(err, &serror) {
					cmd.Printf("Successfull created locally, but got server error: %s. ", serror.Err)
					return
				}
				cmd.PrintErrln(err)
				return
			}

			cmd.Println("Text data created")
		} else {
			cmd.PrintErrln("Not implemented yet")
		}
	},
}

func (dt *dataType) isValid() bool {
	updt := strings.ToUpper(string(*dt))
	if updt == "TEXT" ||
		updt == "CREDENTIALS" ||
		updt == "BANK_CARD" ||
		updt == "BINARY" {
		return true
	}

	return false
}
