package cliclient

import (
	"errors"
	"strings"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common/cliclient"
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

		dt := dataType(strings.ToUpper(datatype))
		if !dt.isValid() {
			cmd.PrintErrln("Invalid data type. Please, use TEXT, CREDENTIALS, BANK_CARD, BINARY.")
			return
		}

		pcd := promptContent{
			errorMsg: "Please, enter description",
			label:    "Enter description",
		}
		description, err := promptTextInput(pcd)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		msg := ""
		if dt == "TEXT" {
			pc := promptContent{
				errorMsg: "Please, enter text",
				label:    "Enter text",
			}
			var text string
			text, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = controller.ucstore.CreateTextData(ctx, text, description)
			msg = "Text data created"
		} else if dt == "BINARY" {
			pc := promptContent{
				errorMsg: "Please, enter file name",
				label:    "Enter file name",
			}
			var filename string
			filename, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = controller.ucstore.CreateBinaryData(ctx, filename, description)
			msg = "Binary data created"
		} else if dt == "CREDENTIALS" {
			pc := promptContent{
				errorMsg: "Please, enter login",
				label:    "Enter login",
			}
			var login, password string
			login, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			pc = promptContent{
				errorMsg: "Please, enter password",
				label:    "Enter password",
			}

			password, err = promptPasswordInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = controller.ucstore.CreateCredentialsData(ctx, clientstore.CredentialDataModel{
				Login:    login,
				Password: password,
			}, description)
			msg = "Credentials created"
		} else if dt == "BANK_CARD" {
			var cardnumber, cvv string
			var pc promptContent
			pc = promptContent{
				errorMsg: "Please, enter card number",
				label:    "Enter card number",
			}

			cardnumber, err = promptTextInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			pc = promptContent{
				errorMsg: "Please, enter cvv",
				label:    "Enter cvv",
			}

			cvv, err = promptPasswordInput(pc)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = controller.ucstore.CreateBankCardData(ctx, clientstore.BankCardDataModel{
				Number: cardnumber,
				CVV:    cvv,
			}, description)
			msg = "Credentials created"
		}

		if err != nil {
			var serror *uccommon.ServerError
			if errors.As(err, &serror) {
				cmd.Printf("Successfull created locally, but got server error: %s. ", serror.Err)
				return
			}
			cmd.PrintErrln(err)
			return
		} else {
			cmd.Println(msg)
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
