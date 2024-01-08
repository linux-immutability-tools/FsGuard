package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/linux-immutability-tools/FsGuard/config"
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
		properties := strings.Split(file, " #FSG# ")

		wg.Add(1)
		go func(prop []string) {
			defer wg.Done()
			if _, err := os.Stat(prop[0]); os.IsNotExist(err) {
				errCh <- fmt.Errorf("[FAIL] %s - File not found", prop[0])
				return
			}

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

			failed := false
			errOut := ""
			if err := validateChecksum(prop[0], sha1sum, prop[1]); err != nil {
				if config.QuitOnFail {
					errCh <- err
					return
				} else {
					failed = true
					errOut = err.Error()
				}
			}

			isSUID, err := strconv.ParseBool(prop[2])
			if err != nil {
				errCh <- fmt.Errorf("[FAIL] %s - Cannot find suid value", prop[0])
				return
			}
			ValidateSUID(prop[0], isSUID)
			if config.QuitOnFail || !failed {
				fmt.Printf("[OK] %s - %s\n", prop[0], sha1sum)
			} else {
				fmt.Printf("%s\n", errOut)
			}
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
