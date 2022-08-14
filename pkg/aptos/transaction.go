package aptos

import (
	"AptosSdk/pkg/rest"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ScriptFunctionPayload https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/ScriptFunctionPayload
type ScriptFunctionPayload struct {
	Type          string   `json:"type"`
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []string `json:"arguments"`
}

// UnsignedTxMsg https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/UserTransactionRequest
type UnsignedTxMsg struct {
	Sender                  string      `json:"sender"`
	SequenceNumber          string      `json:"sequence_number"`
	MaxGasAmount            string      `json:"max_gas_amount"`
	GasUnitPrice            string      `json:"gas_unit_price"`
	GasCurrencyCode         string      `json:"gas_currency_code"`
	ExpirationTimestampSecs string      `json:"expiration_timestamp_secs"`
	Payload                 interface{} `json:"payload"`
}

type SignMessageRsp struct {
	Message string `json:"message"`
}

// TxSignature https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/UserTransactionSignature
type TxSignature struct {
	Type      string `json:"type"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

// TxRequest https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/SubmitTransactionRequest
type TxRequest struct {
	UnsignedTxMsg
	Signature TxSignature `json:"signature"`
}

type GetTxCommonRsp struct {
	Type string `json:"type"`
}

// SubmitTxRsp https://fullnode.devnet.aptoslabs.com/spec.html#/schemas/Transaction
type SubmitTxRsp struct {
	Type                    string `json:"type"`
	Hash                    string `json:"hash"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	GasCurrencyCode         string `json:"gas_currency_code"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	// Payload
	// Signature
}

func Transfer(account *Account, to string, amount int) error {
	account.AccountUpdateCoreResource()

	payload := ScriptFunctionPayload{
		Type:          "script_function_payload",
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments:     []string{to, strconv.Itoa(amount)},
	}
	unSignedTx := UnsignedTxMsg{
		Sender:                  account.PublicAddress(),
		SequenceNumber:          account.CoreResource.SequenceNumber,
		MaxGasAmount:            "1000",
		GasUnitPrice:            "1",
		GasCurrencyCode:         "XUS",
		ExpirationTimestampSecs: strconv.FormatInt(time.Now().Add(600*time.Second).Unix(), 10),
		Payload:                 payload,
	}

	// https://fullnode.devnet.aptoslabs.com/spec.html#/operations/create_signing_message
	url := fmt.Sprintf("%s/transactions/signing_message", FullNodeUrl)
	var signMessage SignMessageRsp
	_, err := rest.DoPost(url, unSignedTx, &signMessage)
	if err != nil {
		fmt.Printf("url [%s] error:%s\n", url, err)
		return err
	}

	toSignMsgBytes, err := hex.DecodeString(signMessage.Message[2:])

	txSignature := TxSignature{
		Type:      "ed25519_signature",
		PublicKey: account.PublicKey(),
		Signature: account.SignMsg(toSignMsgBytes),
	}

	txRequest := TxRequest{
		UnsignedTxMsg: unSignedTx,
		Signature:     txSignature,
	}

	// submit transaction
	// https://fullnode.devnet.aptoslabs.com/spec.html#/operations/submit_transaction
	submitUrl := fmt.Sprintf("%s/transactions", FullNodeUrl)
	var submitTxRsp SubmitTxRsp
	_, err = rest.DoPost(submitUrl, txRequest, &submitTxRsp)
	if err != nil {
		fmt.Printf("submit url [%s] error:%s\n", submitUrl, err)
		return err
	}

	WaitForTransaction(submitTxRsp.Hash)

	return nil
}

// IsTransactionPending https://fullnode.devnet.aptoslabs.com/spec.html#/operations/get_transaction
func IsTransactionPending(txHash string) bool {
	url := fmt.Sprintf("%s/transactions/%s", FullNodeUrl, txHash)
	var rsp GetTxCommonRsp
	statusCode, err := rest.DoGet(url, &rsp)
	if err != nil || statusCode == http.StatusNotFound {
		return true
	}
	return rsp.Type == "pending_transaction"
}

// WaitForTransaction https://fullnode.devnet.aptoslabs.com/spec.html#/operations/get_transaction
func WaitForTransaction(txHash string) bool {
	for isPending, counter := IsTransactionPending(txHash), 0; isPending == true && counter < 10; counter++ {
		time.Sleep(1 * time.Second)
	}

	url := fmt.Sprintf("%s/transactions/%s", FullNodeUrl, txHash)

	var rsp GetTxCommonRsp
	statusCode, err := rest.DoGet(url, &rsp)
	if statusCode != http.StatusOK || err != nil {
		return false
	}

	return true
}
