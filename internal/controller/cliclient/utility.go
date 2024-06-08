package cliclient

import (
	"fmt"

	"github.com/spf13/cobra"
)

func validateRequiredStringFlag(cmd *cobra.Command, flag string) (string, error) {
	f, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}

	if f == "" {
		return "",fmt.Errorf("required flag(s) \"%s\" not set", flag)
	}
	
	return f, nil
}
