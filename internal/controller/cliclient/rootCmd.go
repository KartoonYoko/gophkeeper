package cliclient

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gophkeeper",
	Short: "Gophkeeper is small password keeper",
	Long: `A client to save password and other data.
Complete documentation is available at https://github.com/KartoonYoko/gophkeeper`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}

