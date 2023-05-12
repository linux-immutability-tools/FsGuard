package core

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func ValidatePath(recipePath string) error {
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
		sha256sum := fmt.Sprintf("%x", hashInBytes)
		fmt.Printf("SHA256 hash of file: %s\n", sha256sum)
		fmt.Printf("Wanted SHA256 hash:  %s\n", strings.TrimSpace(properties[1]))
		if strings.Compare(strings.TrimSpace(sha256sum), strings.TrimSpace(properties[1])) == 0 {
			fmt.Println("Checksum Matches!")
		} else {
			fmt.Println("Checksum does not match!")
		}
		err = file.Close()
		if err != nil {
			return err
		}
		fmt.Println()
	}

	return nil
}
