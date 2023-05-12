package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func ValidatePath(recipePath string) error {
	data, err := os.ReadFile(recipePath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(file) == "" {
			continue
		}
		properties := strings.Split(file, " ")

		wg.Add(1)
		go func(prop []string) {
			defer wg.Done()
			file, err := os.Open(prop[0])
			if err != nil {
				errCh <- err
				return
			}
			defer file.Close()

			sha1sum, err := calculateHash(file)
			if err != nil {
				errCh <- err
				return
			}

			prop[1] = strings.TrimSpace(prop[1])

			if err := validateChecksum(prop[0], sha1sum, prop[1]); err != nil {
				errCh <- err
				return
			}

			fmt.Printf("[OK] %s - %s\n", prop[0], sha1sum)
		}(properties)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func calculateHash(file *os.File) (string, error) {
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:20]
	return strings.TrimSpace(fmt.Sprintf("%x", hashInBytes)), nil
}

func validateChecksum(file string, sha1sum, expectedSum string) error {
	if strings.Compare(strings.TrimSpace(sha1sum), expectedSum) != 0 {
		return fmt.Errorf("[FAIL] %s - %s\n\tExpected: %s", file, sha1sum, expectedSum)
	}
	return nil
}
