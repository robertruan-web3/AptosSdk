package aptos

import (
	"AptosSdk/pkg/rest"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strconv"
	"strings"
)

type Account struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey

	CoreResource AccountCoreResource
}

// AccountCoreResource https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/Account
type AccountCoreResource struct {
	SequenceNumber    string `json:"sequence_number"`
	AuthenticationKey string `json:"authentication_key"`
}

func NewAccount(seed string) (*Account, error) {
	seedReader := rand.Reader
	if seed != "" {
		seedReader = strings.NewReader(seed)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(seedReader)
	if err != nil {
		return nil, err
	}

	newAccount := Account{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return &newAccount, nil
}

func (account *Account) PublicKey() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(account.publicKey))
}

func (account *Account) PublicAddress() string {
	hasher := sha3.New256()
	hasher.Write(account.publicKey)
	hasher.Write([]byte("\x00"))
	return fmt.Sprintf("0x%s", hex.EncodeToString(hasher.Sum(nil)))
}

func FoundAccount(account *Account, amount int64) {
	url := fmt.Sprintf("%s/mint?address=%s&amount=%d", NodeUrl, account.PublicAddress(), amount)

	var txHashes []string
	_, err := rest.DoPost(url, txHashes, nil)
	if err != nil {
		fmt.Println("FoundAccount error:", err)
		return
	}

	for _, txHash := range txHashes {
		WaitForTransaction(txHash)
	}

	fmt.Println("FoundAccount successful")
}

func (account *Account) AccountUpdateCoreResource() {
	url := fmt.Sprintf("%s/accounts/%s", FullNodeUrl, account.PublicAddress())
	_, err := rest.DoGet(url, &account.CoreResource)
	if err != nil {
		fmt.Printf("AccountUpdateCoreResource [%s] error:%s \n", url, err)
	}
}

func AccountGetBalance(account *Account) (int64, error) {
	resourceType := "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	url := fmt.Sprintf("%s/accounts/%s/resource/%s", FullNodeUrl, account.PublicAddress(), resourceType)
	var rsp AccountBalanceRsp
	_, err := rest.DoGet(url, &rsp)
	if err != nil {
		return 0, err
	}

	coinValue, _ := strconv.ParseInt(rsp.Data.Coin.Value, 10, 64)
	return coinValue, nil
}

func (account *Account) SignMsg(msg []byte) string {
	signedMsg := ed25519.Sign(account.privateKey, msg)
	return hex.EncodeToString(signedMsg)
}
