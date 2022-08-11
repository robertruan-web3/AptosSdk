package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
	"golang.org/x/crypto/sha3"
	"strings"
)

func GenNewAccount(seed string) (string, string, error) {
	seedReader := rand.Reader
	if seed != "" {
		seedReader = strings.NewReader(seed)
	}

	publicKey, privateKey, err := sign.GenerateKey(seedReader)
	if err != nil {
		fmt.Println("nacl GenerateKey error:", err)
		return "", "", err
	}

	pubKeyBytes := [32]byte(*publicKey)
	priKeyBytes := [64]byte(*privateKey)

	var toHashData []byte
	toHashData = append(toHashData, pubKeyBytes[:]...)
	toHashData = append(toHashData, byte(0x00))

	hasher := sha3.New256()
	hasher.Write(toHashData)

	hashedBytes := hasher.Sum(nil)

	pubKeyHex := hex.EncodeToString(hashedBytes)
	priKeyHex := hex.EncodeToString(priKeyBytes[:])

	return pubKeyHex, priKeyHex, nil
}
