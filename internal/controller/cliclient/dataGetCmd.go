package cliclient

import (
	"github.com/spf13/cobra"
	commondatatype "github.com/KartoonYoko/gophkeeper/internal/common/datatype"
)

func init() {
	dataGetCmd.Flags().String("dataid", "", "Set data id that you want get")

	dataCmd.AddCommand(dataGetCmd)
}

var dataGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get personal data by ID",
	Long:  `Get personal data by ID. Data will be displayed in terminal.`,
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

		// TODO вместо этого метода вызывать метод, который возвращает только тип без самих данных
		item, err := controller.ucstore.GetDataByID(ctx, dataid)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if item.Datatype == commondatatype.DATATYPE_TEXT {
			r, err := controller.ucstore.GetTextDataByID(ctx, dataid)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Description: %s\nData: %s\n", r.Description, r.Text)
		} else if item.Datatype == commondatatype.DATATYPE_BINARY {
			r, err := controller.ucstore.GetBinaryDataByID(ctx, dataid)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Description: %s\nData: %s\n", r.Description, r.Text)
		} else if item.Datatype == commondatatype.DATATYPE_CREDENTIALS {
			r, err := controller.ucstore.GetCredentialsDataByID(ctx, dataid)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Description: %s\nLogin: %s\nPassword: %s\n", r.Description, r.Login, r.Password)
		} else if item.Datatype == commondatatype.DATATYPE_BANK_CARD {
			r, err := controller.ucstore.GetBankCardDataByID(ctx, dataid)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Description: %s\nNumber: %s\nCVV: %s\n", r.Description, r.Number, r.CVV)
		}
	},
}
