package cmd

import (
	"fmt"
	"github.com/linux-immutability-tools/FsGuard/core"
	"github.com/spf13/cobra"
	"os"
)

func NewVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "verify",
		Short:        "Verify the root filesystem based on the given verification file",
		RunE:         ValidateCommand,
		SilenceUsage: true,
	}

	return cmd
}

func ValidateCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no verification file specified")
	}

	fsGuardPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// TODO: Actually verify the list with these values
	signatureFile, err := core.GetSignatureFile(fsGuardPath)
	if err != nil {
		return (err)
	}
	fmt.Println(signatureFile)
	signatureHash, err := core.GetSignatureHash(fsGuardPath)
	if err != nil {
		return err
	}
	fmt.Printf(signatureHash)

	recipePath := args[0]

	err = core.ValidatePath(recipePath)
	if err != nil {
		return err
	}

	return nil
}
