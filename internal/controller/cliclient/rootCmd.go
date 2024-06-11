package cliclient

import (
	"errors"
	"strings"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gophkeeper",
	Short: "Gophkeeper is small password keeper",
	Long: `A client to save password and other data.
Complete documentation is available at https://github.com/KartoonYoko/gophkeeper`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		commandsToSynchronize := []string{
			"data",
		} 

		execute := false
		for _, command := range commandsToSynchronize {
			execute = strings.Contains(cmd.CommandPath(), command)
			if execute {
				break
			}
		}

		if !execute {
			return
		}
		
		cmd.Println("Syncronizing...")
		err := controller.ucstore.Synchronize(cmd.Context())
		if err != nil {
			var serror *clientstore.ServerError
			if errors.As(err, &serror) {
				cmd.Printf("got server error: %v\n\n", serror.Err)
				return
			}

			cmd.Printf("Got error during syncronization: %s\n\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}
