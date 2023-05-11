package cmd

import (
	"crypto/sha256"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
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

	data, err := os.ReadFile(recipePath)
	if err != nil {
		return err
	}

	for _, file := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(file) == "" {
			continue
		}
		properties := strings.Split(file, " ")
		fmt.Println("Path:", properties[0], " Checksum:", properties[1], "is SUID:", properties[2])
		file, err := os.Open(properties[0])
		if err != nil {
			fmt.Println(err)
			return nil
		}

		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			fmt.Println(err)
			return nil
		}

		hashInBytes := hash.Sum(nil)[:32]
		fmt.Printf("SHA256 hash of file: %x\n", hashInBytes)
		fmt.Printf("Wanted hash: 	     %s\n", strings.TrimSpace(properties[1]))
		if strings.Compare(strings.TrimSpace(string(hashInBytes)), strings.TrimSpace(properties[1])) == 0 {
			fmt.Println("Checksum Matches!")
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
