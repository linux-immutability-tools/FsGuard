package core

import (
	"bytes"
	"os"
)

func GetSignatureFile(binary string) (string, error) {
	signatureFileSeperator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}
	signatureFileEndSeperator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x73, 0x65, 0x63, 0x6F, 0x6E, 0x64, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}

	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}

	signatureFileIndex := bytes.Index(data, signatureFileSeperator) + len(signatureFileSeperator)
	signatureFileEndIndex := bytes.Index(data, signatureFileEndSeperator) + len(signatureFileEndSeperator)

	signatureFile := ""

	for i := 0; i < signatureFileEndIndex-signatureFileIndex-len(signatureFileEndSeperator); i++ {
		signatureFile = signatureFile + string(data[signatureFileIndex+i])
	}
	return signatureFile, nil
}

func GetSignatureHash(binary string) (string, error) {
	signatureHashSeperator := []byte{0x2D, 0x2D, 0x2D, 0x2D, 0x62, 0x65, 0x67, 0x69, 0x6E, 0x20, 0x73, 0x65, 0x63, 0x6F, 0x6E, 0x64, 0x20, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x2D, 0x2D, 0x2D, 0x2D}

	data, err := os.ReadFile(binary)
	if err != nil {
		return "", err
	}
	signatureHashIndex := bytes.Index(data, signatureHashSeperator) + len(signatureHashSeperator)

	signatureHash := ""

	for i := 0; i < len(data)-signatureHashIndex; i++ {
		signatureHash = signatureHash + string(data[signatureHashIndex+i])
	}
	return signatureHash, nil

}
