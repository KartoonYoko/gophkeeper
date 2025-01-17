package cliclient

import "github.com/spf13/cobra"

func init() {
	dataCmd.AddCommand(dataGetListCmd)
}

var dataGetListCmd = &cobra.Command{
	Use:   "getlist",
	Short: "Get list of common data (ID, Data type, Description)",
	Long:  `Get list of common data (ID, Data type, Description)`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		// usecase get list
		list, err := controller.ucstore.GetDataList(ctx)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Printf("%-40s %-10s %-10s\n", "ID", "DATATYPE", "DESCRIPTION")
		for _, item := range list {
			cmd.Printf("%-40s %-10s %-10s\n", item.ID, item.Datatype, item.Description)
		}
	},
}
