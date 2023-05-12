package cmd

import (
	"fmt"

	"github.com/linux-immutability-tools/FsGuard/core"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the root filesystem based on the given verification file",
		RunE:  validateCommand,
	}

	return cmd
}

func validateCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no verification file specified")
	}

	recipePath := args[0]

	err := core.ValidatePath(recipePath)
	if err != nil {
		return err
	}

	return nil
}
