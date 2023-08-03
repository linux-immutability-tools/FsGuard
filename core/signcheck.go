package core

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jedisct1/go-minisign"
)

var (
	signatureFileSeparator    = []byte("----begin attach----")
	signatureFileEndSeparator = []byte("----begin second attach----")
	signatureHashSeparator    = []byte("----begin second attach----")
)

func GetSignatureFile(binary string) (string, error) {
	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}

	signatureFileIndex := bytes.LastIndex(data, signatureFileSeparator) + len(signatureFileSeparator)
	signatureFileEndIndex := bytes.LastIndex(data, signatureFileEndSeparator) + len(signatureFileEndSeparator)
	signatureFile := ""

	for i := 0; i < signatureFileEndIndex-signatureFileIndex-len(signatureFileEndSeparator); i++ {
		signatureFile = signatureFile + string(data[signatureFileIndex+i])
	}

	if strings.TrimSpace(signatureFile) == "" {
		fmt.Println("NO SIGNATURE FILE")
		return "", fmt.Errorf("no signature file found")
	}

	return strings.Replace(signatureFile, "----begin second attach---", "", 1), nil
}

func GetSignatureHash(binary string) (string, error) {
	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}

	signatureHashIndex := bytes.LastIndex(data, signatureHashSeparator) + len(signatureHashSeparator)
	signatureHash := ""

	for i := 0; i < len(data)-signatureHashIndex; i++ {
		signatureHash = signatureHash + string(data[signatureHashIndex+i])
	}

	if strings.TrimSpace(signatureHash) == "" {
		fmt.Println("NO SIGNATURE HASH")
		return "", fmt.Errorf("no signature hash found")
	}

	return strings.Replace(signatureHash, "----begin attach---", "", 1), nil
}

func VerifySignature(pubKey string, signature string, filePath string) error {
	pk, err := minisign.NewPublicKey(pubKey)
	if err != nil {
		return err
	}

	sig, err := minisign.DecodeSignature(signature)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	verified, err := pk.Verify(data, sig)
	if err != nil || !verified {
		fmt.Println("Signature verification failed:", err)
		return err
	}

	fmt.Println("Signature verification succeeded!")
	return nil
}
