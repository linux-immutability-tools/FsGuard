package core

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/linux-immutability-tools/FsGuard/config"
)

func validatePathThread(dataCh chan string, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	var file *os.File = nil
	var err error = nil
	for data := range dataCh {
		if len(data) == 0 {
			return
		}

		idx := strings.Index(data, " #FSG# ")
		if idx == -1 {
			errCh <- fmt.Errorf("[FAIL] %s - Malformed line", data)
			continue
		}

		name := data[:idx]
		sig := data[idx+7 : idx+47]
		var isSUID bool
		switch data[idx+54:] {
		case "true":
			isSUID = true
		case "false":
			isSUID = false
		default:
			errCh <- fmt.Errorf("[FAIL] %s - Cannot find suid value", data[idx+54:])
			continue
		}

		if file, err = os.Open(name); err != nil {
			if os.IsNotExist(err) {
				errCh <- fmt.Errorf("[FAIL] %s - File not found", name)
				continue
			}
			errCh <- err
			continue
		}

		sha1sum, err := calculateHash(file)
		if err != nil {
			file.Close()
			errCh <- err
			continue
		}
		file.Close()

		failed := false
		if err = validateChecksum(name, sha1sum, sig); err != nil {
			if config.QuitOnFail {
				errCh <- err
				continue
			} else {
				failed = true
			}
		}

		ValidateSUID(name, isSUID)
		if !failed {
			fmt.Printf("[OK] %s - %s\n", name, sha1sum)
		} else {
			fmt.Printf("%s\n", err)
		}
	}
}

func ValidatePath(recipePath string) error {
	data, err := os.ReadFile(recipePath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	threads := runtime.NumCPU()
	errCh := make(chan error)
	dataCh := make(chan string, threads)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go validatePathThread(dataCh, errCh, &wg)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		dataCh <- scanner.Text()
	}

	for i := 0; i < threads; i++ {
		dataCh <- ""
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(dataCh)
	}()

	for err := range errCh {
		if err != nil {
			fmt.Println(err)
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
