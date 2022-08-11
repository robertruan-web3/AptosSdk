package main

import (
	"AptosSdk/pkg/utils"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
	"golang.org/x/crypto/sha3"
	"strings"
)

func main() {
	seedStr := "abcd1234abcd1234abcd1234abcd1234"

	seed := strings.NewReader(seedStr)
	publicKey, privateKey, err := sign.GenerateKey(seed)
	if err != nil {
		fmt.Println("nacl GenerateKey error:", err)
		return
	}

	pubKeyBytes := [32]byte(*publicKey)
	priKeyBytes := [64]byte(*privateKey)

	var toHashData []byte
	toHashData = append(toHashData, pubKeyBytes[:]...)
	toHashData = append(toHashData, byte(0x00))

	hasher := sha3.New256()
	hasher.Write(toHashData)

	hashedBytes := hasher.Sum(nil)
	fmt.Println("ret:", len(hashedBytes))

	pubKeyHex := hex.EncodeToString(hashedBytes)
	priKeyHex := hex.EncodeToString(priKeyBytes[:])

	fmt.Println(len(pubKeyHex), len(priKeyHex))
	fmt.Println(pubKeyHex)
	fmt.Println(priKeyHex)

	pub2, pri2, err := utils.GenNewAccount(seedStr)
	fmt.Println(pub2)
	fmt.Println(pri2)

	if pubKeyHex != pub2 {
		panic("public key not equal")
	}
	if priKeyHex != pri2 {
		panic("private key not equal")
	}

	fmt.Println("\n=== end of main ===")
}
