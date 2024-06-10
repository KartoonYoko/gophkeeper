package cliclient

import "github.com/spf13/cobra"

func init() {
	dataGetCmd.Flags().String("dataid", "", "Set data id that you want get")

	dataCmd.AddCommand(dataGetCmd)
}

var dataGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Allows to get personal data",
	Long:  `Allows to get personal data`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := cmd.Context()

		dataid, err := cmd.Flags().GetString("dataid")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if dataid == "" {
			cmd.Help()
			return
		}

		item, err := controller.ucstore.GetDataByID(ctx, dataid)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Printf("%-40s %-10s %-10s %-10s\n", "ID", "DATATYPE", "DATA", "DESCRIPTION")
		cmd.Printf("%-40s %-10s %-10s %-10s\n", item.ID, item.Datatype, item.Data, item.Description)			
	},
}
