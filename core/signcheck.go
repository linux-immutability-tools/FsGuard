package core

import (
	"bytes"
	"fmt"
	"github.com/jedisct1/go-minisign"
	"os"
	"strings"
)

func GetSignatureFile(binary string) (string, error) {
	signatureFileSeparator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}
	signatureFileEndSeparator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x73, 0x65, 0x63, 0x6F, 0x6E, 0x64, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}

	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}
	signatureFileIndex := bytes.Index(data, signatureFileSeparator) + len(signatureFileSeparator)
	signatureFileEndIndex := bytes.Index(data, signatureFileEndSeparator) + len(signatureFileEndSeparator)
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
	signatureHashSeparator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x73, 0x65, 0x63, 0x6F, 0x6E, 0x64, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}

	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}
	signatureHashIndex := bytes.Index(data, signatureHashSeparator) + len(signatureHashSeparator)

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
