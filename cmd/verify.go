package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the root filesystem based on the given verification file",
		RunE:  validateCommand,
	}

	return cmd
}

type File struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
	IsSUID   bool   `json:"isSUID"`
}

func validateCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no verification file specified")
	}

	recipePath := args[0]

	data, err := os.ReadFile(recipePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	var files []File
	err = json.Unmarshal(data, &files)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}
	for _, file := range files {
		fmt.Println("Path: ", file.Path, " Checksum: ", file.Checksum, " is SUID: ", file.IsSUID)
	}
	//fmt.Println("Name:", file[0].Path)
	//fmt.Println("Age:", person.Age)

	return nil
}
