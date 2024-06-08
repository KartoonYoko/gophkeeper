package cliclient

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	flagRegisterLogin    string
	flagRegisterPassword string
)

type promptContent struct {
	errorMsg string
	label    string
}

func init() {
	registerCmd.PersistentFlags().StringVar(&flagRegisterLogin, "login", "", "set your login to authenticate")
	registerCmd.PersistentFlags().StringVar(&flagRegisterPassword, "password", "", "set your password to authenticate")

	root.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var flagRegisterLogin, flagRegisterPassword string
		ctx := cmd.Context()

		err = cmd.MarkFlagRequired("login")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		err = cmd.ValidateRequiredFlags()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if flagRegisterPassword == "" {
			wordPromptContent := promptContent{
				"Please provide a password.",
				"Enter a password:",
			}
			flagRegisterPassword, err = promptGetInput(wordPromptContent)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		err = controller.ucauth.Register(ctx, flagRegisterLogin, flagRegisterPassword)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("register successfully")
	},
}

func promptGetInput(pc promptContent) (string, error) {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
		Mask:      '*',
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	fmt.Println("")

	return result, nil
}
