package cliclient

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	dataCreateCmd.Flags().String("datatype", "TEXT", "Set data type that you want create. Possible values: TEXT, CREDENTIALS, BANK_CARD, BINARY.")
	dataCmd.AddCommand(dataCreateCmd)
}

type dataType string

var dataCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// ctx := cmd.Context()

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

		// todo usecase create
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